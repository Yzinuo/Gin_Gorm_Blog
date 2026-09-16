## 博客设计

响应式：移动优先，大屏适应

## 黑红主题与 3D 自我介绍

- 首页首屏、导航与作者卡片均可进入 `/about`，原成就滚动区已移除。
- `/about` 使用独立布局，支持滚动镜头、眼神跟随、八个经历节点及返回博客。后台 `/front/about` 的文字介绍仍显示在页尾。
- 主题变量集中在 `src/styles/theme.css`；UnoCSS 的 `brand`、`surface`、`foreground` 等颜色引用同一组变量。
- 经历文案在 `src/views/about/content.js`。其顺序需与 `camera-map.json`、模型中的 `CameraAction` 动画保持一致。
- 模型和贴纸位于 `public/resume/`，源自提供的 viewer。模型约 28 MB，只在进入关于页时请求；无需外部 3D CDN。
- 场景加载失败或设备不支持 WebGL 时，经历文字仍可阅读，并提供重试。离开页面会取消下载、停止动画并释放图形资源。

本地运行：`npm run dev`；生产检查：`npm run build`。配置的本地后端地址为 `http://localhost:8765`，未启动时首页会显示文章重试入口。

## TODO

- 不使用 NaiveUI 作为组件库，使用自研组件库 ✅

阅读README.md: 通常项目会有一个README文件，里面包含了项目的介绍、安装步骤和使用说明。

查看main.js: 这是Vue应用的入口文件，可以了解Vue实例的初始化过程，以及引入了哪些插件和组件。

查看router.js: 了解应用的路由结构，每个路由对应哪个页面组件。

查看store/index.js: 了解Vuex的初始化过程，以及状态管理的模块划分。

查看components/和views/: 了解项目的组件结构，特别是页面级别的组件。

查看api.js: 了解与后端API交互的逻辑。
