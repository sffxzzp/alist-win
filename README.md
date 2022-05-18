alist-win
======
一个基于 [go-webview2](https://github.com/jchv/go-webview2) 的 [alist](https://github.com/Xhofe/alist) 构建

会自动检测新版并自动构建，版本检测采用 [PipeDream](https://pipedream.com)

请自行前往 [Actions](https://github.com/sffxzzp/alist-win/actions) 中的最新构建中下载 artifact

提示
------
1. 会在启动后创建 `password.txt`，包含后台密码
2. 会在 `data` 目录下创建 `EBWebView` 文件夹，用于存储 WebView2 的数据
3. 修改了 upx 的压缩参数，进一步缩减文件体积
4. 其他和正常使用无太大区别
