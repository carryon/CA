package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/bocheninc/CA/deploy/components/utils"
	"github.com/bocheninc/CA/deploy/tables"
	"github.com/bocheninc/L0/components/crypto"
)

var (
	defaultKeyFileName  = "ca.key"
	defaultCertFileName = "ca.crt"
)

type Ca struct {
	db              *sql.DB
	dirPath         string
	keyPath         string
	certPath        string
	rootPrivateKey  *rsa.PrivateKey
	rootCertificate *x509.Certificate
}

func NewCa(db *sql.DB, path string) *Ca {
	return &Ca{db: db,
		dirPath:  path,
		keyPath:  filepath.Join(path, defaultKeyFileName),
		certPath: filepath.Join(path, defaultCertFileName)}
}

func (c *Ca) Init() error {
	_, err := utils.OpenDir(c.dirPath)
	if err != nil {
		return fmt.Errorf("open dir %s error: %s", c.dirPath, err)
	}

	//file is not exist generate new cert and key of read local file
	if !utils.FileExist(c.keyPath) || !utils.FileExist(c.certPath) {
		certificate := crypto.NewCertificate(crypto.CertInformation{IsCA: true})

		//creat key
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("create key error: %s", err)
		}

		//create cert
		cert, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &rsaKey.PublicKey, rsaKey)
		if err != nil {
			return fmt.Errorf("create root cert error: %s", err)
		}
		if err := write(c.certPath, "CERTIFICATE", cert); err != nil {
			return fmt.Errorf("write root cert err: %s", err)
		}

		key := x509.MarshalPKCS1PrivateKey(rsaKey)
		if err := write(c.keyPath, "PRIVATE KEY", key); err != nil {
			return fmt.Errorf("write root key err: %s", err)
		}

		certificate, err = x509.ParseCertificate(cert)
		if err != nil {
			return fmt.Errorf("parse root cert  error: %s", err)
		}
		c.rootCertificate = certificate
		c.rootPrivateKey = rsaKey
	} else {
		keyBytes, err := ioutil.ReadFile(c.keyPath)
		if err != nil {
			return fmt.Errorf("read key file error: %s", err)
		}

		certBytes, err := ioutil.ReadFile(c.certPath)
		if err != nil {
			return fmt.Errorf("read cert file error: %s", err)
		}

		rsaKey, err := crypto.ParseKey(keyBytes)
		if err != nil {
			return fmt.Errorf("prase key err: %s", err)
		}

		certificate, err := crypto.ParseCrt(certBytes)
		if err != nil {
			return fmt.Errorf("prase cert err: %s", err)

		}

		c.rootCertificate = certificate
		c.rootPrivateKey = rsaKey
	}

	//clear t_cert table
	tx, _ := c.db.Begin()
	cert := new(tables.Cert)
	if err := cert.DeleteAll(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("init clear table Cert error: %s", err)
	}
	tx.Commit()

	return nil

}

func (c *Ca) GetCert(chainID, nodeID, pubKey string) ([]byte, error) {
	certs, err := tables.QueryCertByID(c.db, chainID, nodeID)
	if err != nil {
		log.Errorf("query Cert by %s:%s error: %s", chainID, nodeID, err)
		return nil, err
	}

	switch len(certs) {
	case 0:
		cert, err := c.generateNodeCert(chainID, nodeID, pubKey)
		if err != nil {
			return nil, err
		}

		certInfo := &tables.Cert{
			PublicKey: pubKey,
			ChainID:   chainID,
			NodeID:    nodeID,
			Crt:       string(cert),
			Created:   time.Now(),
		}

		tx, _ := c.db.Begin()
		if err = certInfo.Insert(tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()

		return cert, nil
	case 1:
		if certs[0].PublicKey == pubKey {
			return []byte(certs[0].Crt), nil
		}

		cert, err := c.generateNodeCert(chainID, nodeID, pubKey)
		if err != nil {
			return nil, err
		}
		certInfo := &tables.Cert{
			PublicKey: pubKey,
			ChainID:   chainID,
			NodeID:    nodeID,
			Crt:       string(cert),
			Created:   time.Now(),
		}

		tx, _ := c.db.Begin()
		if err = certInfo.UpdateCert(tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()

		return cert, nil

	default:
		return nil, fmt.Errorf("query Cert by %s:%s not only one ", chainID, nodeID)
	}
}

func (c *Ca) generateNodeCert(chainID, nodeID, pubKey string) ([]byte, error) {
	pblock, _ := pem.Decode([]byte(pubKey))
	publicKey, err := x509.ParsePKIXPublicKey(pblock.Bytes)
	if err != nil {
		log.Errorf("parse publicKey %s ,err: %s", pubKey, err)
		return nil, err
	}

	log.Debugf("create cert by %s:%s ", chainID, nodeID)
	baseInfo := crypto.CertInformation{Locality: []string{chainID, nodeID}}
	certificate := crypto.NewCertificate(baseInfo)

	crt, err := x509.CreateCertificate(rand.Reader, certificate, c.rootCertificate, publicKey.(*rsa.PublicKey), c.rootPrivateKey)
	if err != nil {
		log.Error("create crt error: ", err)
		return nil, err
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crt,
	}

	buffer := new(bytes.Buffer)

	pem.Encode(buffer, block)

	return buffer.Bytes(), nil
}

func (c *Ca) GetRootCertificate() ([]byte, error) {
	certBytes, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return nil, fmt.Errorf("read cert file error: %s", err)
	}

	return certBytes, nil
}

func write(filename, Type string, data []byte) error {
	if utils.FileExist(filename) {
		if err := os.RemoveAll(filename); err != nil {
			return err
		}
	}
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	return pem.Encode(File, &pem.Block{Bytes: data, Type: Type})
}
