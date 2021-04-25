package com.java.dao;

import com.java.model.CasbinRule;
import org.apache.ibatis.annotations.*;
import org.springframework.stereotype.Repository;

import java.util.Collection;
import java.util.List;
@Mapper
@Repository // 这个注解可以不加, 纯粹是为了 idea 在 autowired 中不报错
public interface CasbinDao {
    // todo  测试自动生成 id
    @Insert("INSERT INTO casbin_rule(id, p_type, v0, v1) VALUES(#{id}, #{pType}, #{v0}, #{v1})")
    @Options(useGeneratedKeys = true, keyProperty = "id", keyColumn = "id")
    int insert(CasbinRule casbin);

    @Update(value = {
            "<script>",
            "UPDATE users",
            "<set>",
            "<if test='username != null'>, username = #{username}</if>",
            "<if test='password != null'>, password = #{password}</if>",
            "</set>",
            "WHERE id = #{id}",
            "</script>"
    })
    int updateById(CasbinRule user);

    @Delete("DELETE FROM casbin_rule WHERE id = #{id}")
    int deleteById(@Param("id") Integer id);

    @Select(value = {
            "<script>",
            "SELECT * FROM casbin_rule",
            "WHERE id IN",
            "<foreach item='id' collection='ids' separator=',' open='(' close=')' index=''>",
            "#{id}",
            "</foreach>",
            "</script>"
    })
    List<CasbinRule> selectByIds(@Param("ids") Collection<Integer> ids);
}
