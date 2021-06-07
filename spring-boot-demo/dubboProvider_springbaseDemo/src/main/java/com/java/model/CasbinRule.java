package com.java.model;

public class CasbinRule {
    private int id;
    private String pType;
    private String v0;
    private String v1;

    public CasbinRule(String pType, String v0, String v1) {
        this.pType = pType;
        this.v0 = v0;
        this.v1 = v1;
    }

    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public String getpType() {
        return pType;
    }

    public void setpType(String pType) {
        this.pType = pType;
    }

    public String getV0() {
        return v0;
    }

    public void setV0(String v0) {
        this.v0 = v0;
    }

    public String getV1() {
        return v1;
    }

    public void setV1(String v1) {
        this.v1 = v1;
    }

    @Override
    public String toString() {
        return "CasbinRule{" +
                "id=" + id +
                ", pType='" + pType + '\'' +
                ", v0='" + v0 + '\'' +
                ", v1='" + v1 + '\'' +
                '}';
    }
}
