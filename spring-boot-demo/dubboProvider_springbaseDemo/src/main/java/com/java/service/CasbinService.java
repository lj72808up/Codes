package com.java.service;

import com.java.dao.CasbinDao;
import com.java.model.CasbinRule;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Collection;
import java.util.List;

@Service
public class CasbinService {
    @Autowired
    private CasbinDao casbinDao;

    @Transactional
    public int insert(CasbinRule casbin) throws InterruptedException {
        int newId = this.casbinDao.insert(casbin);
        System.out.println("插入的 newId: " + newId);
        Thread.sleep(10 * 1000);
        System.out.println("事务已提交");
        return newId;
    }

    public int updateById(CasbinRule casbin) {
        return this.casbinDao.updateById(casbin);
    }

    public List<CasbinRule> selectByIds(Collection<Integer> ids) {
        return this.casbinDao.selectByIds(ids);
    }


    public CasbinDao getCasbinDao() {
        return casbinDao;
    }

    public void setCasbinDao(CasbinDao casbinDao) {
        this.casbinDao = casbinDao;
    }
}
