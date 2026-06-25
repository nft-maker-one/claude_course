from openai import OpenAI

# 在此处硬编码填写你的 API Key
# 示例: API_KEY = "sk-proj-xxxx..."
API_KEY = "sk-proj-123456"

# 使用提供的 API Key 初始化客户端
client = OpenAI(api_key=API_KEY)

def send_hello_world():
    try:
        # 调用 Chat Completions API
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",  # 你也可以将其替换为 "gpt-4" 或 "gpt-4o"
            messages=[
                {"role": "user", "content": "hello world"}
            ]
        )
        
        # 提取并打印模型回复的内容
        reply = response.choices[0].message.content
        print("OpenAI 回复:")
        print(reply)
        
    except Exception as e:
        print(f"调用 API 时发生错误: {e}")

if __name__ == "__main__":
    if not API_KEY:
        print("提示: 请先在代码中填入你的 API Key。")
    else:
        send_hello_world()