<div align="center">

# 碎碎 SuiSui

**轻量级备忘录应用 · Markdown · 标签 · 图片上传 · 消息反应 · 热力图**

![Vue 3](https://img.shields.io/badge/Vue_3-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vuetify 4](https://img.shields.io/badge/Vuetify_4-1867C0?style=flat-square&logo=vuetify&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)

</div>

---

## Features

| 功能 | 说明 |
|------|------|
| Markdown 编辑 | 粗体、斜体、标题、代码块、链接、列表、引用 |
| 多图上传 | 多图上传 + 轮播查看 + 点击放大 |
| 消息反应 | emoji 库，游客可用 |
| 标签分类 | 标签筛选 + 全文搜索高亮 |
| 置顶排序 | 重要笔记置顶 |
| 活动热力图 | 每日备忘数量可视化 |
| 回收站 | 软删除 + 恢复 + 永久删除 |
| 用户系统 | 注册 / 登录 / 角色权限 |
| 自定义主题色 | 每位用户独立设置 |
| 暗色模式 | 手动切换 |
| 移动适配 | 响应式设计 |
| 后台管理 | 系统设置、用户管理 |
| 单文件部署 | 一个二进制即可运行 |

---

## Quick Start

### Development

```bash
npm install
cd server-go && go run main.go
npx vite --port 5173 --host
```

Open **http://localhost:5173** - Default admin: `admin / admin`

### Production Build

```bash
npm run build
cd server-go && go build -o suisui .
./suisui              # default port 3001
./suisui -port 8080   # custom port
```

---

## Project Structure

```
src/
  stores/        # Pinia stores (auth, notes, settings)
  views/         # Pages (NotesPage, AdminPage)
  components/    # Components (NoteCard, Heatmap, MarkdownPreview...)
server-go/
  main.go        # Go backend (single file)
  dist/          # Built frontend assets
```

---

## Tech Stack

| Frontend | Backend |
|----------|---------|
| Vue 3 + TypeScript | Go |
| Vuetify 4 | SQLite |
| Pinia | RESTful API |
| Marked + Highlight.js | Embedded frontend in binary |
| Vite 6 | Zero external dependencies |

---

## TODO

- [ ] Infinite scroll / pagination
- [ ] Auto-save drafts
- [ ] Code block theme switching
- [x] Note export / import
- [ ] Drag-and-drop pin ordering
- [ ] WebSocket real-time updates

---

<div align="center">

**碎碎** - Capture every spark of inspiration

MIT License

</div>
