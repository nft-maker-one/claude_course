# ============================================================
# 演示文件 2: Layer 1 — Python 危险 API 检测
# 预期触发规则: pickle_deserialization, unsafe_yaml_load,
#              python_subprocess_shell, os_system_injection
# ============================================================
import pickle
import yaml
import subprocess
import os


def load_user_data(filepath):
    """从文件加载用户上传的数据 — 反序列化 RCE"""
    with open(filepath, 'rb') as f:
        return pickle.load(f)


def parse_config(raw_yaml):
    """解析 YAML 配置 — 可通过 !!python/object 执行任意代码"""
    return yaml.load(raw_yaml)


def search_files(user_query):
    """搜索文件 — shell=True 命令注入"""
    result = subprocess.run(f"grep -r '{user_query}' /data", shell=True, capture_output=True)
    return result.stdout.decode()


def ping_host(host):
    """Ping 主机 — os.system 命令注入"""
    os.system(f"ping -c 1 {host}")
