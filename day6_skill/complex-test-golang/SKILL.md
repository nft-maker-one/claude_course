---
name: complex-test-golang
description: 自动运行 Go 项目的单元测试 + 性能基准测试，并在当前目录生成 Markdown 格式的完整性能报告。
disable-model-invocation: true
---

你是一个 Go 测试自动化专家。请按照以下步骤在当前目录运行 Go 测试并生成报告。

## 第一步：确认环境

```bash
!`go version`
!`ls go.mod 2>/dev/null && echo "go.mod found" || echo "ERROR: no go.mod, wrong directory?"`
```

如果没有 `go.mod`，停止并提示用户切换到正确目录。

## 第二步：运行单元测试

```bash
!`go test -v -count=1 ./... 2>&1 | tee /tmp/go_unit_out.txt; echo "EXIT:$?"`
```

统计通过 / 失败用例数，记录结果。

## 第三步：运行性能基准测试

```bash
!`go test -bench=. -benchmem -benchtime=3s -count=1 ./... 2>&1 | tee /tmp/go_bench_out.txt`
```

## 第四步：生成 Markdown 报告

将以下信息写入 `bench_report.md`（用 Write 工具）：

1. **环境信息表格**：Go 版本、OS/Arch、CPU、时间戳、单元测试总状态（PASS / FAIL）
2. **单元测试结果**：只列出每个子测试的 PASS/FAIL 一行摘要，失败的用 ❌ 标注
3. **基准测试分析**，按算法分节：
   - AES-GCM Encrypt / Decrypt（128 / 192 / 256 位 × 消息大小）
   - ECDSA Sign / Verify（P-256 / P-384 / P-521）
   - RSA Encrypt / Decrypt / Sign / Verify（1024 / 2048 / 4096 位）
4. **性能对比表格**：每个算法族的相对速度（以最快的为基准 = 1×）
5. **关键洞察**：3-5 条简短结论（例如：AES key-size 对 GCM 吞吐影响 < 5%；RSA-4096 解密比 RSA-2048 慢 7.6×）
6. **完整原始输出**（折叠在 `<details>` 标签内）

报告完成后输出：`报告已写入 bench_report.md，共 X 个 benchmark，单元测试：PASS/FAIL`
