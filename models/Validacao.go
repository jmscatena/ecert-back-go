package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type CertVal struct {
	ID            uint64      `gorm:"primary_key;auto_increment" json:"id"`
	CertificadoID uint64      `json:"certificado_id"`
	Certificado   Certificado `json:"certificado"`
	Validado      time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"validado_em"`
	Hash          string      `gorm:"size:100;not null;" json:"chave"`
}

func (v *CertVal) Validate() error {

	if v.CertificadoID == 0 {
		return errors.New("obrigatório: idcertificado")
	}

	return nil
}

func (v *CertVal) Create(db *gorm.DB) (int64, error) {
	var verr error
	if verr = v.Validate(); verr != nil {
		return -1, verr
	}
	err := db.Debug().Model(&CertVal{}).Create(&v).Error
	if err != nil {
		return 0, err
	}
	return int64(v.ID), nil
}

func (v *CertVal) List(db *gorm.DB) (*[]CertVal, error) {
	var err error
	CertVals := []CertVal{}
	err = db.Debug().Model(&CertVal{}).Limit(100).Find(&CertVals).Error
	if err != nil {
		return &[]CertVal{}, err
	}
	return &CertVals, nil
}

func (v *CertVal) Find(db *gorm.DB, uid uint64) (*CertVal, error) {
	err := db.Debug().Model(&CertVal{}).Where("id = ?", uid).Take(&v).Error
	if err != nil {
		return &CertVal{}, err
	}
	return v, nil
}

func (v *CertVal) Update(db *gorm.DB, uid uint64) (*CertVal, error) {
	err := db.Debug().Model(&CertVal{}).Where("id = ?", uid).Take(&CertVal{}).UpdateColumns(
		map[string]interface{}{
			"CertificadoID": v.CertificadoID,
			"Certificado":   v.Certificado,
		},
	).Error
	if err != nil {
		return &CertVal{}, err
	}
	return v, nil
}

func (v *CertVal) Delete(db *gorm.DB, pid uint64, uid uint64) (int64, error) {

	db = db.Debug().Model(&CertVal{}).Where("id = ? and apresentador_id = ?", pid, uid).Take(&CertVal{}).Delete(&CertVal{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
