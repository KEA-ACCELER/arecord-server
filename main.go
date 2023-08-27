package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/KEA-ACCELER/arecord-server/pkg/elastic"
	"github.com/KEA-ACCELER/arecord-server/pkg/types/record"
	"github.com/KEA-ACCELER/arecord-server/pkg/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go"
	qp "github.com/quic-s/quics-protocol"
	pb "github.com/quic-s/quics-protocol/proto/v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server, err := qp.New()
	if err != nil {
		log.Panicln(err)
	}

	current, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("current directory: ", current)

	// 파일 데이터를 저장할 디렉토리 초기화
	dataDir := filepath.Join(current, "data")
	log.Println("data directory: ", dataDir)

	_, err = os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dataDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created data folder")
	} else {
		log.Println("Using existing data folder")
	}

	err = server.RecvMessage(func(conn quic.Connection, message *pb.Message) {
		log.Println("message received")
		log.Println(message.Type, string(message.Message))
		switch message.Type {
		case "removed":
			// dotlog 파일ㅇ르 읽어서 timestamp를 가져온다.
			dotlog, err := os.ReadFile(filepath.Join(dataDir, string(message.Message)+".log"))
			if err != nil {
				log.Panicln(err)
			}
			timestamp, err := strconv.Atoi(string(dotlog))
			if err != nil {
				log.Panicln(err)
			}
			timestamp++
			err = os.WriteFile(filepath.Join(dataDir, string(message.Message)+".log"), []byte(strconv.Itoa(timestamp)), 0755)
			if err != nil {
				log.Panicln(err)
			}
			err = os.WriteFile(filepath.Join(dataDir, string(message.Message)+"."+strconv.Itoa(timestamp)), nil, 0755)
			if err != nil {
				log.Panicln(err)
			}

			diff, insert, delete := utils.Diff(
				filepath.Join(dataDir, string(message.Message)+"."+strconv.Itoa(timestamp-1)),
				filepath.Join(dataDir, string(message.Message)+"."+strconv.Itoa(timestamp)),
			)

			// get hash
			hasher := sha256.New()
			hasher.Write(nil)
			hash := hasher.Sum(nil)
			hashString := hex.EncodeToString(hash)

			// create record and send to elastic search
			record := &record.Record{
				Hash:      hashString,
				Version:   timestamp,
				Path:      string(message.Message),
				Diff:      diff,
				Time:      time.Now(),
				Editor:    strings.Split(conn.RemoteAddr().String(), ":")[0],
				Size:      0,
				Insert:    insert,
				Delete:    delete,
				Extension: filepath.Ext(string(message.Message)),
			}

			u, err := uuid.NewRandom()
			if err != nil {
				log.Panicln(err)
			}

			elastic.CreateDoc(u.String(), record.MarshalJson())
		case "giveme":
			log.Println("giveme")
		}

	})
	if err != nil {
		log.Panicln(err)
	}

	err = server.RecvFile(func(conn quic.Connection, fileInfo *pb.FileInfo, fileBuf []byte) {
		log.Println("file received: ", fileInfo.Path)
		log.Println("data: ", fileBuf)

		// 파일 생성인 경우, 파일 수정인 경우
		_, err = os.Stat(filepath.Join(dataDir, fileInfo.Path+".log"))
		if os.IsNotExist(err) {
			dir, _ := path.Split(fileInfo.Path)
			err = os.MkdirAll(filepath.Join(dataDir, dir), 0755)
			if err != nil {
				log.Panicln(err)
			}

			err = os.WriteFile(filepath.Join(dataDir, fileInfo.Path+".log"), []byte("0"), 0755)
			if err != nil {
				log.Panicln(err)
			}
			err = os.WriteFile(filepath.Join(dataDir, fileInfo.Path+".0"), nil, 0755)
			if err != nil {
				log.Panicln(err)
			}
		}

		// dotlog 파일ㅇ르 읽어서 timestamp를 가져온다.
		dotlog, err := os.ReadFile(filepath.Join(dataDir, fileInfo.Path+".log"))
		if err != nil {
			log.Panicln(err)
		}
		timestamp, err := strconv.Atoi(string(dotlog))
		if err != nil {
			log.Panicln(err)
		}
		timestamp++
		err = os.WriteFile(filepath.Join(dataDir, fileInfo.Path+".log"), []byte(strconv.Itoa(timestamp)), 0755)
		if err != nil {
			log.Panicln(err)
		}
		err = os.WriteFile(filepath.Join(dataDir, fileInfo.Path+"."+strconv.Itoa(timestamp)), fileBuf, 0755)
		if err != nil {
			log.Panicln(err)
		}

		diff, insert, delete := utils.Diff(
			filepath.Join(dataDir, fileInfo.Path+"."+strconv.Itoa(timestamp-1)),
			filepath.Join(dataDir, fileInfo.Path+"."+strconv.Itoa(timestamp)),
		)

		// get hash
		hasher := sha256.New()
		hasher.Write(fileBuf)
		hash := hasher.Sum(nil)
		hashString := hex.EncodeToString(hash)

		// create record and send to elastic search
		record := &record.Record{
			Hash:      hashString,
			Version:   timestamp,
			Path:      fileInfo.Path,
			Diff:      diff,
			Time:      time.Now(),
			Editor:    strings.Split(conn.RemoteAddr().String(), ":")[0],
			Size:      uint64(len(fileBuf)),
			Insert:    insert,
			Delete:    delete,
			Extension: filepath.Ext(fileInfo.Path),
		}

		u, err := uuid.NewRandom()
		if err != nil {
			log.Panicln(err)
		}

		elastic.CreateDoc(u.String(), record.MarshalJson())
	})
	if err != nil {
		log.Panicln(err)
	}

	// start server
	err = server.Listen("0.0.0.0", 18080)
	if err != nil {
		log.Panicln(err)
	}
}
