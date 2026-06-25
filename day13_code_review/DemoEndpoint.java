// ============================================================
// 演示文件 4: Layer 2/3 — Java SSRF + SQL注入 (正则无法覆盖)
// Layer 1 无法检测这些 Java 特有漏洞模式
// 需要 Layer 2 (LLM Diff Review) 和 Layer 3 (Agentic Review) 补充
// ============================================================


import java.io.*;
import java.net.*;
import java.sql.*;


public class DemoEndpoint {

    // 漏洞 1: SSRF — 用户输入直接作为 URL 发起请求
    // 攻击者可以访问内网: ?url=http://169.254.169.254/latest/meta-data/
    public String fetchUrl(String url) throws Exception {
        URL u = new URL(url);
        HttpURLConnection conn = (HttpURLConnection) u.openConnection();
        BufferedReader reader = new BufferedReader(
            new InputStreamReader(conn.getInputStream()));
        StringBuilder sb = new StringBuilder();
        String line;
        while ((line = reader.readLine()) != null) {
            sb.append(line);
        }
        return sb.toString();
    }

    // 漏洞 2: SQL 注入 — 字符串拼接构造 SQL
    // 攻击者: ?id=1 OR 1=1; DROP TABLE users;--
    public String getUser(String id) throws Exception {
        Connection conn = DriverManager.getConnection(
            "jdbc:mysql://localhost/mydb", "root", "s3cretPassw0rd!");
        Statement stmt = conn.createStatement();
        String sql = "SELECT * FROM users WHERE id = " + id;
        ResultSet rs = stmt.executeQuery(sql);
        StringBuilder sb = new StringBuilder();
        while (rs.next()) {
            sb.append(rs.getString("username")).append("\n");
        }
        conn.close();
        return sb.toString();
    }

    // 漏洞 3: 路径遍历 — 未校验的文件读取
    // 攻击者: ?name=../../../etc/passwd
    public String readFile(String name) throws Exception {
        File f = new File("/var/uploads/" + name);
        return new String(java.nio.file.Files.readAllBytes(f.toPath()));
    }
}
