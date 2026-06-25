// ============================================================
// 演示文件 5: Layer 3 — 跨文件漏洞 (Agentic Review 专长)
// 命令执行 + 不安全反序列化
// Agentic Reviewer 可以 Read/Grep 追踪数据流
// ============================================================


import java.io.*;


public class DemoAdvanced {

    // 漏洞 4: 命令执行 — 用户输入直接传入 shell
    // 攻击者: ?cmd=id;cat /etc/shadow
    public String execCommand(String cmd) throws Exception {
        Process proc = Runtime.getRuntime().exec(
            new String[]{"/bin/sh", "-c", cmd});
        BufferedReader reader = new BufferedReader(
            new InputStreamReader(proc.getInputStream()));
        StringBuilder sb = new StringBuilder();
        String line;
        while ((line = reader.readLine()) != null) {
            sb.append(line).append("\n");
        }
        return sb.toString();
    }

    // 漏洞 5: 不安全反序列化 — Java ObjectInputStream
    // 攻击者可构造恶意序列化对象实现 RCE
    public String deserialize(byte[] data) throws Exception {
        ByteArrayInputStream bais = new ByteArrayInputStream(data);
        ObjectInputStream ois = new ObjectInputStream(bais);
        Object obj = ois.readObject();
        ois.close();
        return obj.toString();
    }

    // 漏洞 6: 不安全的文件上传路径
    // 攻击者: filename=../../webapps/ROOT/shell.jsp
    public String upload( String filename,
                         byte[] content) throws Exception {
        String uploadDir = "/var/uploads/";
        File dest = new File(uploadDir + filename);
        FileOutputStream fos = new FileOutputStream(dest);
        fos.write(content);
        fos.close();
        return "Saved to: " + dest.getAbsolutePath();
    }
}
