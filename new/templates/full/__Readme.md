{{ .ProName }}系统
================

[![pipeline status]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/badges/master/pipeline.svg)]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/commits/master)
[![coverage report]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/badges/master/coverage.svg)]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/commits/master)

### 前端地址说明
| 环境         | 地址                                                  |
| ------------ | ----------------------------------------------------- |
| QA           | https://{{.ProNameUrlSafe}}.qa.{{ .Domain }}          |
| QA2          | https://{{.ProNameUrlSafe}}.qa2.{{ .Domain }}         |
| Production   | https://{{.ProNameUrlSafe}}.{{ .Domain }}             |


### 后端地址说明
| 环境         | 地址                                                         |
| ------------ | ------------------------------------------------------------ |
| QA           | https://{{.ProNameUrlSafe}}.api-qa.{{ .Domain }}             |
| QA2          | https://{{.ProNameUrlSafe}}.api-qa2.{{ .Domain }}            |
| Production   | https://{{.ProNameUrlSafe}}.api.{{ .Domain }}                |


# 模块1接口列表
- [样例接口](routes/example-api-doc.md)

# 模块2接口列表
- [样例接口](routes/example-api-doc.md)
