package surfstore

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) {
	// synchronize local
	createIndexFile(client.BaseDir)
	localMetaMap, err := syncLocal(client.BaseDir, client.BlockSize)
	if err != nil {
		log.Fatal(err)
	}
	localMetaMap = syncRemote(localMetaMap, client)
	err = WriteMetaFile(localMetaMap, client.BaseDir)
	if err != nil {
		log.Fatal(err)
	}
}

func createIndexFile(dir string) {
	metaFilePath, err := filepath.Abs(ConcatPath(dir, DEFAULT_META_FILENAME))
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(metaFilePath)
	if err != nil {
		f, errCreate := os.Create(metaFilePath)
		if errCreate != nil {
			log.Fatal(errCreate)
		}
		errClose := f.Close()
		if errClose != nil {
			log.Fatal(errClose)
		}
	}
}

func syncLocal(dir string, bufSize int) (fileMetaMap map[string]*FileMetaData, e error) {
	absDir, err := filepath.Abs(dir)
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		log.Fatal(err)
	}
	oldLocalIndex, err := LoadMetaFromMetaFile(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileInDir := make(map[string]bool, 0)
	for _, file := range files {
		if file.Name() == DEFAULT_META_FILENAME {
			continue
		}
		fileInDir[file.Name()] = true
		fileHashList := make([]string, 0)
		absPath, err := filepath.Abs(ConcatPath(dir, file.Name()))
		f, err := os.Open(absPath)
		if err != nil {
			log.Fatal(err)
		}
		for {
			buffer := make([]byte, bufSize)
			n, err := f.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			hashString := GetBlockHashString(buffer[:n])
			fileHashList = append(fileHashList, hashString)
		}
		metaData, ok := oldLocalIndex[file.Name()]
		if ok {
			if same := compareHashList(fileHashList, metaData.BlockHashList); !same {
				metaData.BlockHashList = fileHashList
				metaData.Version += 1
				oldLocalIndex[file.Name()] = metaData
			}
		} else {
			//if len(fileHashList) == 0 {
			//	oldLocalIndex[file.Name()] = &FileMetaData{Filename: file.Name(), Version: 1, BlockHashList: []string{""}}
			//} else {
			//	oldLocalIndex[file.Name()] = &FileMetaData{Filename: file.Name(), Version: 1, BlockHashList: fileHashList}
			//}
			oldLocalIndex[file.Name()] = &FileMetaData{Filename: file.Name(), Version: 1, BlockHashList: fileHashList}
		}
	}
	for oldFile, oldFileData := range oldLocalIndex {
		if _, ok := fileInDir[oldFile]; !ok {
			//oldFileData := oldLocalIndex[oldFile]
			if !compareHashList(oldFileData.BlockHashList, []string{"0"}) {
				oldFileData.Version += 1
				oldFileData.BlockHashList = []string{"0"}
			}
		}
	}
	return oldLocalIndex, nil
}

func compareHashList(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, s := range l1 {
		if s != l2[i] {
			return false
		}
	}
	return true
}

func syncRemote(localMap map[string]*FileMetaData, client RPCClient) map[string]*FileMetaData {
	var remoteMap map[string]*FileMetaData
	err := client.GetFileInfoMap(&remoteMap)
	if err != nil {
		log.Fatal(err)
	}

	for localFile, localData := range localMap {
		_, ok := remoteMap[localFile]
		if (!ok) || (ok && localData.Version == remoteMap[localFile].Version+1) {
			upLoadBlock(client, localData)
			var ver int32
			err := client.UpdateFile(localData, &ver)
			if err != nil {
				log.Fatal(err)
			}
			if ver == -1 {
				var remoteMap map[string]*FileMetaData
				err = client.GetFileInfoMap(&remoteMap)
				if err != nil {
					log.Fatal(err)
				}
				localMap[localFile] = remoteMap[localFile]
				writeLocalFile(client, localMap[localFile])
			} else if ver != localData.Version {
				log.Fatal("Unexpected Wrong Version")
			}
		} else {
			localMap[localFile] = remoteMap[localFile]
			writeLocalFile(client, localMap[localFile])
		}
	}
	for remoteFile, remoteData := range remoteMap {
		if _, ok := localMap[remoteFile]; !ok {
			localMap[remoteFile] = remoteData
			writeLocalFile(client, remoteData)
		}
	}
	return localMap
}

func writeLocalFile(client RPCClient, localData *FileMetaData) {
	localPath, err := filepath.Abs(ConcatPath(client.BaseDir, localData.Filename))
	if err != nil {
		log.Fatal(err)
	}

	if (len(localData.BlockHashList) == 1) && (localData.BlockHashList[0] == "0") {
		if _, err := os.Stat(localPath); err != nil {
			return
		}
		err = os.Remove(localPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	fileBytes := make([]byte, 0)
	var blockStoreAddr string
	if err := client.GetBlockStoreAddr(&blockStoreAddr); err != nil {
		log.Fatal(err)
	}

	for _, hash := range localData.BlockHashList {
		var block Block
		if err := client.GetBlock(hash, blockStoreAddr, &block); err != nil {
			log.Fatal(err)
		}
		fileBytes = append(fileBytes, block.BlockData[:block.BlockSize]...)
	}

	err = os.WriteFile(localPath, fileBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func upLoadBlock(client RPCClient, localData *FileMetaData) {
	if (len(localData.BlockHashList) == 1) && (localData.BlockHashList[0] == "0") {
		return
	}
	var blockAddr string
	var succ bool
	err := client.GetBlockStoreAddr(&blockAddr)
	if err != nil {
		log.Fatal(err)
	}
	absPath, err := filepath.Abs(ConcatPath(client.BaseDir, localData.Filename))
	f, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	hashBlockMap := make(map[string][]byte)
	hashs := make([]string, 0)
	for {
		buffer := make([]byte, client.BlockSize)
		n, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		hash := GetBlockHashString(buffer[:n])
		hashBlockMap[hash] = buffer[:n]
		hashs = append(hashs, hash)

	}
	existBlock := make([]string, 0)
	err = client.HasBlocks(hashs, blockAddr, &existBlock)
	if err != nil {
		log.Fatal(err)
	}
	existBlockBoolMap := make(map[string]bool, 0)
	for _, b := range existBlock {
		existBlockBoolMap[b] = true
	}
	for h, b := range hashBlockMap {
		if _, ok := existBlockBoolMap[h]; ok {
			continue
		}
		err = client.PutBlock(&Block{BlockData: b, BlockSize: int32(len(b))}, blockAddr, &succ)
		if err != nil {
			log.Fatal(err)
		}
	}

}
