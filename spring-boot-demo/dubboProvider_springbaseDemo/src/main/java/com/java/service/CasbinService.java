package com.java.service;

import com.java.dao.CasbinDao;
import com.java.model.CasbinRule;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Collection;
import java.util.List;

@Service
public class CasbinService {
    @Autowired
    private CasbinDao casbinDao;

    public int insert(CasbinRule casbin){
        return this.casbinDao.insert(casbin);
    }
    public int updateById(CasbinRule casbin){
        return this.casbinDao.updateById(casbin);
    }
    public List<CasbinRule> selectByIds(Collection<Integer> ids){
        return this.casbinDao.selectByIds(ids);
    }


    public CasbinDao getCasbinDao() {
        return casbinDao;
    }

    public void setCasbinDao(CasbinDao casbinDao) {
        this.casbinDao = casbinDao;
    }
}
