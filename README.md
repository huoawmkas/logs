#简单logs日志包 使用
######dir: ../log         # 文件保存路径
######file: logs          # 文件名称,实际会保存为{filename}+{datetime}
######level: 3            # 日志等级：0-error，1-warning，2-info，3-debug
######savefile: false     # 是否保存为文件，置为false会输出到标准输出

logs.Init(dir, file, level, savefile)

logs.Error(err1,err2,err3)