package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Certificado struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	EventoID       uint64    `gorm:"not null" json:"evento_id"`
	Evento         Evento    `json:"evento"`
	ParticipanteID uint64    `gorm:"not null" json:"participante_id"`
	Participante   Usuario   `json:"participante"`
	Validacao      string    `json:"validacao"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Certificado) Validate() error {

	if p.EventoID < 1 {
		return errors.New("obrigatório: evento")
	}
	if p.ParticipanteID < 1 {
		return errors.New("obrigatório: apresentador")
	}
	return nil
}

func (p *Certificado) Create(db *gorm.DB) (int64, error) {
	if verr := p.Validate(); verr != nil {
		return -1, verr
	}
	err := db.Debug().Model(&Certificado{}).Create(&p).Error
	if err != nil {
		return 0, err
	}
	return int64(p.ID), nil
}

func (p *Certificado) Update(db *gorm.DB, uid uint64) (*Certificado, error) {
	err := db.Debug().Model(&Certificado{}).Where("id = ?", uid).Take(&Certificado{}).UpdateColumns(
		map[string]interface{}{
			"EventoID":       p.EventoID,
			"Evento":         p.Evento,
			"ParticipanteID": p.ParticipanteID,
			"Participante":   p.Participante,
			"UpdatedAt":      time.Now()}).Error

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Certificado) List(db *gorm.DB) (*[]Certificado, error) {
	Certificados := []Certificado{}
	//err := db.Debug().Model(&Certificado{}).Limit(100).Find(&Certificados).Error
	result := db.Find(&Certificados)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Certificados, nil
}

func (p *Certificado) Find(db *gorm.DB, uid uint64) (*Certificado, error) {
	err := db.Debug().Model(&Certificado{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Certificado{}, err
	}
	return p, nil
}

func (p *Certificado) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Certificado, error) {
	Certificados := []Certificado{}
	params := strings.Split(param, ";")
	uids := uid[0].([]interface{})
	if len(params) != len(uids) {
		return nil, errors.New("condição inválida")
	}
	result := db.Where(strings.Join(params, " AND "), uids...).Find(&Certificados)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Certificados, nil
}

func (p *Certificado) Delete(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Delete(&Certificado{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Certificado) DeleteBy(db *gorm.DB, cond string, uid uint64) (int64, error) {
	result := db.Delete(&Certificado{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
