package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Evento struct {
	ID             uint64      `gorm:"primary_key;auto_increment" json:"id"`
	Descricao      string      `gorm:"size:500;not null" json:"descricao"`
	Apresentador   Usuario     `json:"apresentador"`
	ApresentadorID uint64      `gorm:"not null" json:"apresentador_id"`
	Local          string      `json:"local"`
	Data           time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"data_evento"`
	Instituicao    Instituicao `json:"instituicao"`
	InstituicaoID  uint64      `gorm:"not null" json:"instituicao_id"`
	CreatedAt      time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Evento) Prepare() {
	p.ID = 0
	p.Local = html.EscapeString(strings.TrimSpace(p.Local))
	p.Apresentador = Usuario{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Evento) Validate() error {

	if p.Descricao == "" {
		return errors.New("obrigatório: evento")
	}
	if p.Local == "" {
		return errors.New("obrigatório: local")
	}
	if p.ApresentadorID < 1 {
		return errors.New("obrigatório: apresentador")
	}
	return nil
}

func (p *Evento) Create(db *gorm.DB) (int64, error) {
	var verr error
	if verr = p.Validate(); verr != nil {
		return -1, verr
	}
	err := db.Debug().Model(&Evento{}).Create(&p).Error
	if err != nil {
		return 0, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.ApresentadorID).Take(&p.Apresentador).Error
		if err != nil {
			return -2, err
		}
	}
	return int64(p.ID), nil
}

func (p *Evento) Update(db *gorm.DB, uid uint64) (*Evento, error) {
	err := db.Debug().Model(&Evento{}).Where("id = ?", uid).Take(&Evento{}).UpdateColumns(
		map[string]interface{}{
			"Descricao":      p.Descricao,
			"Apresentador":   p.Apresentador,
			"ApresentadorID": p.ApresentadorID,
			"Local":          p.Local,
			"Data":           p.Data,
			"Instituicao":    p.Instituicao,
			"InstituicaoID":  p.InstituicaoID,
			"UpdatedAt":      time.Now(),
		}).Error
	if err != nil {
		return &Evento{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.ApresentadorID).Take(&p.Apresentador).Error
		if err != nil {
			return &Evento{}, err
		}
	}
	return p, nil
}

func (p *Evento) List(db *gorm.DB) (*[]Evento, error) {
	Eventos := []Evento{}
	err := db.Debug().Model(&Evento{}).Limit(100).Find(&Eventos).Error
	if err != nil {
		return &[]Evento{}, err
	}
	if len(Eventos) > 0 {
		for i := range Eventos {
			err := db.Debug().Model(&Usuario{}).Where("id = ?", Eventos[i].ApresentadorID).Take(&Eventos[i].Apresentador).Error
			if err != nil {
				return &[]Evento{}, err
			}
		}
	}
	return &Eventos, nil
}

func (p *Evento) Find(db *gorm.DB, pid uint64) (*Evento, error) {
	err := db.Debug().Model(&Evento{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Evento{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.ApresentadorID).Take(&p.Apresentador).Error
		if err != nil {
			return &Evento{}, err
		}
	}
	return p, nil
}

func (p *Evento) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Evento, error) {
	Eventos := []Evento{}
	params := strings.Split(param, ";")
	uids := uid[0].([]interface{})
	if len(params) != len(uids) {
		return nil, errors.New("condição inválida")
	}
	result := db.Where(strings.Join(params, " AND "), uids...).Find(&Eventos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Eventos, nil
}

func (p *Evento) Delete(db *gorm.DB, uid uint64) (int64, error) {

	db = db.Debug().Model(&Evento{}).Where("id = ? ", uid).Take(&Evento{}).Delete(&Evento{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Evento) DeleteBy(db *gorm.DB, cond string, uid uint64) (int64, error) {
	result := db.Delete(&Evento{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
