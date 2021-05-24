#### 统计页面pv
1.  计算 pv 的
    ```bash
    # 计算 pv 的
    cat AStar.info.2020-12-15.001.log | grep BaseController.go:22 | grep -v "/jones/redash/timestamp" | awk '{ split($9,a,"/") ; count[a[2]]++;} END {for (i in count) {print i "===>" count[i]}}'
    ```
    pv : 2384
2. 计算 uv 的
    ```bash
    # 计算 uv 的
    cat AStar.info.2020-12-15.001.log | grep BaseController.go:22 | grep -v "/jones/redash/timestamp" | awk '{ count[$5]++;} END {for (i in count) {print i "===>" count[i]}}' | wc -l 
    ```
    uv : 34 个
