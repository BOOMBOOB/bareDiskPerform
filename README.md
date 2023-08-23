# bareDiskPerform
## 项目介绍
用于服务器裸盘性能测试，并将数据收集到mysql数据库中

考虑到某些场景下SAS通道带宽可能存在瓶颈的问题，目前仅支持单任务顺序执行



## 代码结构

- main目录：代码运行入口main.go及配置文件config.json所在
- module目录：日志、fio、数据库、等模块代码所在，由于结构比较简单所以没有细分




## 使用方法
1. 将代码拉取到Linux服务器上，进入main.go所在目录，进行静态编译

   ```
   go build -tags 'mysql' -o disktest main.go
   ```

2. 将编译出的二进制文件disktest和config.json放到测试服务器上同一目录下

3. 对应修改config.json配置文件

4. 运行测试

   ```
   ./disktest
   ```



## 配置说明

- level：配置日志打印等级，可配置值为："debug", "info", "warn", "error", "dpanic", "panic", "fatal"
- mysql：配置存储测试信息
    - server：mysql所在服务器IP地址
    - port：mysql使用端口
    - database：数据库名称
    - username：用户名
    - password：密码
- disks：配置测试磁盘信息
    - mode：测试模式，可配置值为："auto"，"manual"。"auto"代表自动模式，会自动扫描测试服务器上所有HDD硬盘，此模式为高危模式，请确保使用此模式时所有HDD上的数据均不需要；"manual"代表手动模式，只测试配置中指定的盘符
    - type：测试类型列表，可配置值为：["read", "write", "randread", "randwrite"]，”read“和”write“模式会自动使用1M的IO块大小，”randread“、”randwrite“模式会自动使用4KB的IO块大小
    - devices：指定测试盘符列表，可配置值为：["sdd", "sdf", ......]
- ramp_time：测试预热时间
- runtime：测试运行时间
- iodepth：测试IO队列深度