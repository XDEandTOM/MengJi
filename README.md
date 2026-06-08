# MengJi

一个轻量级备忘录应用，支持 Markdown、标签、图片上传、热力图、暗色模式等特性。

## 技术栈

| 前端 | 后端 |
|------|------|
| Vue 3 + TypeScript | Go |
| Vuetify 4 | SQLite |
| Pinia | 单文件编译 |
| Marked + Highlight.js | RESTful API |
| Vite 6 | |

## 功能

- 备忘录 Markdown 编辑与渲染
- 多图上传 + Carousel 轮播 + 点击放大
- 标签分类与搜索
- 置顶与删除
- 活动热力图
- 用户注册/登录
- 后台管理
- 暗色/亮色主题切换
- 移动端适配
- 单二进制部署

## 快速开始

### 开发模式

```bash
npm install
cd server-go && go run main.go
npx vite --port 5173 --host
```

打开 http://localhost:5173

### 生产构建

```bash
npm run build
cd server-go && go build -o mengji .
./mengji
```

默认管理员：admin / admin

## License

MIT
