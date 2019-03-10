package core

import (
	"bytes"
	"encoding/hex"

	"github.com/danitello/go-blockchain/common/byteutil"
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/core/types"

	"github.com/dgraph-io/badger"
)

var (
	utxoPrefix = []byte("utxo-")
)

/** utxo_set is additional database functions for BlockChain involving the running collection of current utxos */

/*Reindex deletes the current UTXOSet and establishes a new one */
func (bc *BlockChain) Reindex() {
	bc.DeleteWithKeyPrefix(utxoPrefix)

	err := bc.ChainDB.Database.Update(func(txn *badger.Txn) error {
		for txID, txos := range bc.GetUTXO() {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			key = append(utxoPrefix, key...)

			err = txn.Set(key, byteutil.Serialize(txos))
			errutil.Handle(err)
		}

		return nil
	})
	errutil.Handle(err)
}

/*DeleteWithKeyPrefix deletes all data whose key is prefixed by a given value
@param prefix - to delete
*/
func (bc *BlockChain) DeleteWithKeyPrefix(prefix []byte) {
	deleteKeys := func(keysToDelete [][]byte) error {
		if err := bc.ChainDB.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysToDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000 // badgerdb
	bc.ChainDB.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysToDelete := make([][]byte, 0, collectSize)
		numKeysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysToDelete = append(keysToDelete, key)
			numKeysCollected++
			if numKeysCollected == collectSize {
				err := deleteKeys(keysToDelete)
				errutil.Handle(err)
				keysToDelete = make([][]byte, 0, collectSize)
				numKeysCollected = 0
			}
		}
		if numKeysCollected > 0 {
			err := deleteKeys(keysToDelete)
			errutil.Handle(err)
		}
		return nil
	})
}

/*UpdateUTXOSet manages adding and deleting tx references in set resulting from new Block
@param block - containing new Transactions to traverse
*/
func (bc *BlockChain) UpdateUTXOSet(block *types.Block) {
	err := bc.ChainDB.Database.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, txin := range tx.Inputs {
					updatedTXO := types.TxOutputs{}
					dbID := append(utxoPrefix, txin.TxID...)
					item, err := txn.Get(dbID)
					errutil.Handle(err)
					v, err := item.Value()
					errutil.Handle(err)

					TXO := types.DeserializeTxOutputs(v)

					for txoIdx, txo := range TXO.Outputs {
						if txoIdx != txin.OutputIdx {
							updatedTXO.Outputs = append(updatedTXO.Outputs, txo)
						}
					}

					if len(updatedTXO.Outputs) == 0 {
						err := txn.Delete(dbID) // No more UTXO
						errutil.Handle(err)
					} else {
						err := txn.Set(dbID, byteutil.Serialize(updatedTXO))
						errutil.Handle(err)
					}
				}
			}
			newTXO := types.TxOutputs{}
			for _, txo := range tx.Outputs {
				newTXO.Outputs = append(newTXO.Outputs, txo) // Just go ahead and add them
			}

			dbID := append(utxoPrefix, tx.ID...)
			err := txn.Set(dbID, byteutil.Serialize(newTXO))
			errutil.Handle(err)
		}

		return nil
	})
	errutil.Handle(err)
	bc.Reindex()
}

/*GetUTXOWithPubKey gets utxos owned by a pub key hash with a total balance up to a given amount
@param pubKeyHash - the pub key hash in question
@param max - the max value to search up to
@return map - the utxos
@return int - the balance
*/
func (bc *BlockChain) GetUTXOWithPubKey(pubKeyHash []byte, max int) (map[string][]int, int) {
	UTXO := make(map[string][]int)
	balance := 0

	err := bc.ChainDB.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			errutil.Handle(err)

			k = bytes.TrimPrefix(k, utxoPrefix)
			txID := hex.EncodeToString(k)
			TXO := types.DeserializeTxOutputs(v)

			for txoIdx, txo := range TXO.Outputs {
				if txo.IsLockedWithKey(pubKeyHash) && balance < max {
					balance += txo.Amount
					UTXO[txID] = append(UTXO[txID], txoIdx)
				}
			}
		}
		return nil
	})
	errutil.Handle(err)

	return UTXO, balance
}

/*CountUTX gets the number of Transactions with UTXO in them
@return the number
*/
func (bc BlockChain) CountUTX() int {
	count := 0

	err := bc.ChainDB.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			count++
		}

		return nil
	})

	errutil.Handle(err)

	return count
}
