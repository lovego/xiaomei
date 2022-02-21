{{ .ProName }}系统
==================

[![pipeline status]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/badges/master/pipeline.svg)]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/commits/master)
[![coverage report]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/badges/master/coverage.svg)]({{ .RepoPrefix }}/{{.ProNameUrlSafe}}/commits/master)

### 地址说明
| 环境         | 地址                                 |
| ------------ | ------------------------------------ |
| QA           | https://qa-{{ .Domain }}             |
| Preview      | https://preview-{{ .Domain }}        |
| Production   | https://{{ .Domain }}                |

### API文档
[API文档](./doc/api)
