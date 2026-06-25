// ============================================================
// 演示文件 1: Layer 1 — JavaScript XSS 模式检测
// 预期触发规则: innerHTML_xss, document_write_xss, eval_injection
// ============================================================

function renderUserComment(comment) {
    // 危险: 直接将用户输入塞进 innerHTML → XSS
    document.getElementById("comments").innerHTML = comment;
}

function legacyRender(html) {
    // 危险: document.write 可被利用进行 XSS
    document.write("<div>" + html + "</div>");
}

function runCalculator(expr) {
    // 危险: eval 执行任意代码
    return eval(expr);
}
