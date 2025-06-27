// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	sumamodels "smlm_automation/pkg/models/susemanager"
)

// GetSlaves - get list of SUSE Manager secondary servers
//
// param: sessionKey
// return:
func (p *Proxy) GetSlaves(sessionKey string) ([]sumamodels.Slaves, error) {
	p.logger.Debug("GetSlaves call started")
	path := "sync/slave/getSlaves"
	response, err := p.suse.SuseManagerCall(nil, "GET", p.cfg.Host, path, sessionKey)
	if err != nil {
		p.logger.Error("error while fetching suse slaves", zap.Any("error", err))
		return nil, fmt.Errorf("error while fetching suse slaves: %s", err)
	}
	var resultSuc []sumamodels.Slaves
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	}
	return resultSuc, nil
}

// SyncSlaveGetSlaveByName - return SUSE Manager Secondary information
//
// param: auth
// param: slaveFQDN
// return:
func (p *Proxy) SyncSlaveGetSlaveByName(auth AuthParams, slaveFQDN string) (sumamodels.Slaves, error) {
	p.logger.Debug("GetSlaves call started")
	var resultSuc sumamodels.Slaves
	path := "sync/slave/getSlaveByName"
	body, _ := json.Marshal(map[string]interface{}{"slaveFqdn": slaveFQDN})
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while fetching suse slave ID", zap.Any("error", err))
		return resultSuc, fmt.Errorf("error while fetching suse slave ID: %s", err)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return resultSuc, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return resultSuc, fmt.Errorf("unable to process the received data. err: %s", err)
		}
		return resultSuc, nil
	}
	return resultSuc, fmt.Errorf("error while fetching SUSE slave ID. Invalid Response code: %d", response.StatusCode)
}

// SyncSlaveDelete - Delete SUSE Manager Secondary in Intersync configuration
//
// param: auth
// param: slaveID
// return:
func (p *Proxy) SyncSlaveDelete(auth AuthParams, slaveID int) (int, error) {
	p.logger.Debug("GetSlaves call started")
	body, _ := json.Marshal(map[string]interface{}{"slaveId": slaveID})
	path := "sync/slave/delete"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while deleting SUSE Manager Slave entry", zap.Any("error", err))
		return 0, fmt.Errorf("error while deleting SUSE Manager Slave entry. err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		p.logger.Error("error while deleting SUSE Manager Slave entry", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("deleting SUSE Manager Slave entry Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncSlaveDelete"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// SyncSlaveCreate - create SUSE Manager Secondary in intersync
//
// param: auth
// param: slaveFQDN
// param: isEnabled
// param: allowAllOrgs
// return:
func (p *Proxy) SyncSlaveCreate(auth AuthParams, slaveFQDN string, isEnabled bool, allowAllOrgs bool) (sumamodels.Slaves, error) {
	p.logger.Debug("SyncSlaveCreate call started")
	var resultSuc sumamodels.Slaves
	path := "sync/slave/create"
	body, _ := json.Marshal(map[string]interface{}{"slaveFqdn": slaveFQDN, "isEnabled": isEnabled, "allowAllOrgs": allowAllOrgs})
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while creating suse slaves", zap.Any("error", err))
		return resultSuc, fmt.Errorf("error while creating suse slaves: %s", err)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return resultSuc, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return resultSuc, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		p.logger.Error("error while creating SUSE Manager Slave entry", zap.Any("StatusCode", response.StatusCode))
		return resultSuc, fmt.Errorf("creating SUSE Manager Slave entry Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncSlaveCreate"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// SyncMasterGetMasterByLabel - get SUSE Manager Primary information
//
// param: auth
// param: label
// return:
func (p *Proxy) SyncMasterGetMasterByLabel(auth AuthParams, label string) (sumamodels.SlavesIssMaster, error) {
	p.logger.Debug("SyncMasterGetMasterByLabel call started")
	var resultSuc sumamodels.SlavesIssMaster
	path := "sync/master/getMasterByLabel"
	body, _ := json.Marshal(map[string]interface{}{"label": label})
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while fetching suse master ID", zap.Any("error", err))
		return resultSuc, fmt.Errorf("error while fetching suse slave ID: %s", err)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return resultSuc, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return resultSuc, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		p.logger.Error("No master defined", zap.Any("StatusCode", response.StatusCode))
		return resultSuc, fmt.Errorf("no Master Defined. So no error")
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncMasterGetMasterByLabel"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// SyncMasterDelete - remove SUSE Manager Primary from intersync
//
// param: auth
// param: masterID
// return:
func (p *Proxy) SyncMasterDelete(auth AuthParams, masterID int) (int, error) {
	p.logger.Debug("SyncMasterDelete call started")
	body, _ := json.Marshal(map[string]interface{}{"masterId": masterID})
	path := "sync/master/delete"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while deleting SUSE Manager master entry", zap.Any("error", err))
		return 0, fmt.Errorf("error while deleting SUSE Manager master entry. err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		p.logger.Debug("calling sync/master/delete Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling sync/master/delete Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncSlaveMaster"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// SyncMasterCreate - add SUSE Manager Primary to intersync
//
// param: auth
// param: label
// return:
func (p *Proxy) SyncMasterCreate(auth AuthParams, label string) (sumamodels.SlavesIssMaster, error) {
	p.logger.Debug("SyncMasterCreate call started")
	var resultSuc sumamodels.SlavesIssMaster
	path := "sync/master/create"
	body, _ := json.Marshal(map[string]interface{}{"label": label})
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while create SUSE Manager master record", zap.Any("error", err))
		return resultSuc, fmt.Errorf("error while create SUSE Manager master record: %s", err)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return resultSuc, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return resultSuc, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		p.logger.Debug("calling sync/master/create Failed", zap.Any("StatusCode", response.StatusCode))
		return resultSuc, fmt.Errorf("calling sync/master/create Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncMasterCreate"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// SyncMasterMakeDefault - make SUSE Manager Primary the default
//
// param: auth
// param: masterID
// return:
func (p *Proxy) SyncMasterMakeDefault(auth AuthParams, masterID int) (int, error) {
	p.logger.Debug("SyncMasterMakeDefault call started")
	body, _ := json.Marshal(map[string]interface{}{"masterId": masterID})
	path := "sync/master/makeDefault"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while setting master SUSE Manager", zap.Any("error", err))
		return 0, fmt.Errorf("error while setting master SUSE Manager: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		p.logger.Debug("SyncMasterMakeDefault call Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling SyncMasterMakeDefault Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "SyncMasterMakeDefault"), zap.Any("response", resultSuc))
	return 1, nil
}

// SyncMasterSetCaCert - set CA certificate from SUSE Manager Primary to intersync
//
// param: auth
// param: masterID
// param: caCert
// return:
func (p *Proxy) SyncMasterSetCaCert(auth AuthParams, masterID int, caCert string) (int, error) {
	p.logger.Debug("SyncMasterSetCaCert call started")
	body, _ := json.Marshal(map[string]interface{}{"masterId": masterID, "caCertFilename": caCert})
	path := "sync/master/setCaCert"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("error while setting master SUSE Manager CaCert", zap.Any("error", err))
		return 0, fmt.Errorf("error while setting master SUSE Manager CaCert: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		p.logger.Debug("Setting CA Cert has failed. HTTP Statuscode not 200", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("setting CA Cert has failed. HTTP Statuscode: %s", fmt.Sprint(response.StatusCode))
	}
	p.logger.Debug("Response from api", zap.Any("api", "sync/master/setCaCert"), zap.Any("response", resultSuc))
	return resultSuc, nil
}
