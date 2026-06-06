import os
import csv
from collections import Counter

def scan_extensions(directory):
    """
    扫描指定目录下的文件并统计各后缀名的数量。
    """
    # TODO: 初始化一个 Counter 对象，用于高效统计各个后缀名的出现次数
    ext_counts = Counter()

    # TODO: 使用 os.walk 遍历目标目录及其所有的子目录
    for root, _, files in os.walk(directory):
        # TODO: 遍历当前目录层级下的每一个文件
        for file in files:
            # TODO: 使用 os.path.splitext 分离出文件名和后缀名 (例如: '.txt')
            _, ext = os.path.splitext(file)
            
            # TODO: 将后缀名统一转换为小写，避免大小写导致重复统计 (例如 '.JPG' 和 '.jpg')
            ext = ext.lower()
            
            # TODO: 判断文件是否存在后缀名，如果为空，则归类为 '无后缀名'
            if not ext:
                ext = '无后缀名'
                
            # TODO: 更新该后缀名在 Counter 中的计数
            ext_counts[ext] += 1

    # TODO: 返回统计完成的数据字典
    return ext_counts

def save_to_csv(data, output_filename):
    """
    将统计数据写入 CSV 文件。
    """
    # TODO: 以写入模式 ('w') 打开目标 CSV 文件，设置 newline='' 以避免 Windows 下出现多余空行，指定 utf-8 编码
    with open(output_filename, mode='w', newline='', encoding='utf-8') as csvfile:
        # TODO: 初始化 CSV 写入器对象
        writer = csv.writer(csvfile)
        
        # TODO: 写入 CSV 文件的表头 (Header)
        writer.writerow(['文件后缀', '数量'])
        
        # TODO: 遍历统计数据字典，将每个后缀及其对应的数量作为一行写入文件
        for ext, count in data.items():
            # TODO: 执行单行写入操作
            writer.writerow([ext, count])

def main():
    """
    主控函数。
    """
    # TODO: 定义需要扫描的目录路径 ('.' 代表当前执行脚本所在的文件夹)
    target_dir = '.'
    
    # TODO: 定义输出结果的 CSV 文件名称
    output_file = 'extension_counts.csv'

    print(f"开始扫描目录: '{os.path.abspath(target_dir)}' ...")
    
    # TODO: 调用 scan_extensions 函数，获取当前目录及子目录下的后缀统计结果
    counts = scan_extensions(target_dir)

    # TODO: 检查是否扫描到了文件，避免空跑
    if not counts:
        print("未在当前文件夹找到任何文件。")
        return

    print(f"扫描完成，开始将结果写入 '{output_file}' ...")
    
    # TODO: 调用 save_to_csv 函数，将统计结果持久化到磁盘中
    save_to_csv(counts, output_file)

    # TODO: 打印最终的成功提示信息
    print("操作成功完成！您可以查看生成的 CSV 文件。")

if __name__ == '__main__':
    # TODO: 验证当前模块是否作为主程序运行，若是则触发 main() 函数
    main()