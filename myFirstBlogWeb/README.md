# 博客系统后端接口文档

## 目录

- 1. 用户相关
  - 1.1 用户注册
  - 1.2 用户登录
  - 1.3 获取验证码
- 2. 文章相关
  - 2.1 创建文章
  - 2.2 获取文章列表
  - 2.3 修改文章
  - 2.4 删除文章
- 3. 评论相关
  - 3.1 获取评论列表
  - 3.2 发表评论
- 4. 通用说明

---

## 1. 用户相关

### 1.1 用户注册

- **接口**：`POST /register`
- **请求体（JSON）**：

  ```json
  {
    "username": "string",        // 用户名
    "password_hash": "string",   // 密码哈希（前端加密后传输）
    "captcha": "string",         // 验证码答案
    "captcha_id": "string"       // 验证码ID
  }
  ```

- **响应**：

  - 成功

    ```json
    {
      "message": "注册成功",
      "userID": "string"
    }
    ```

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

### 1.2 用户登录

- **接口**：`POST /login`
- **请求体（JSON）**：

  ```json
  {
    "username": "string",
    "password_hash": "string",
    "captcha": "string",
    "captcha_id": "string"
  }
  ```

- **响应**：

  - 成功

    ```json
    {
      "token": "jwt-token-string",
      "userID": "string",
      "username": "string"
    }
    ```

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

### 1.3 获取验证码

- **接口**：`POST /captcha`
- **响应**：

  ```json
  {
    "captcha_id": "string",
    "image_base64": "data:image/png;base64,..."
  }
  ```

---

## 2. 文章相关

### 2.1 创建文章

- **接口**：`POST /ctx/articles`
- **请求头**：`Authorization: Bearer <token>`
- **请求体（JSON）**：

  ```json
  {
    "title": "string",
    "content": "string"
  }
  ```

- **响应**：

  - 成功

    ```json
    {
      "article_id": "string",
      "title": "string",
      "content": "string",
      "author_id": "string",
      "create_time": "2025-08-10T12:00:00Z",
      "modify_time": "2025-08-10T12:00:00Z"
    }
    ```

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

### 2.2 获取文章列表

- **接口**：`GET /ctx/articles?page=1&limit=10`
- **响应**：

  ```json
  {
    "articles": [
      {
        "article_id": "string",
        "title": "string",
        "content": "string",
        "author_id": "string",
        "create_time": "2025-08-10T12:00:00Z",
        "modify_time": "2025-08-10T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 10
  }
  ```

---

### 2.3 修改文章

- **接口**：`PUT /ctx/articles/{id}`
- **请求头**：`Authorization: Bearer <token>`
- **请求体（JSON）**：

  ```json
  {
    "title": "string",
    "content": "string"
  }
  ```

- **响应**：

  - 成功：HTTP 204 No Content

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

### 2.4 删除文章

- **接口**：`DELETE /ctx/articles/{id}`
- **请求头**：`Authorization: Bearer <token>`
- **响应**：

  - 成功：HTTP 204 No Content

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

## 3. 评论相关

### 3.1 获取评论列表

- **接口**：`GET /ctx/comments/{articleId}?page=1&limit=10`
- **响应**：

  ```json
  {
    "comments": [
      {
        "comment_id": "string",
        "content": "string",
        "author_id": "string",
        "article_id": "string",
        "create_time": "2025-08-10T12:00:00Z",
        "modify_time": "2025-08-10T12:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "limit": 10
  }
  ```

---

### 3.2 发表评论

- **接口**：`POST /ctx/comments/{articleId}`
- **请求头**：`Authorization: Bearer <token>`
- **请求体（JSON）**：

  ```json
  {
    "article_id": "string",
    "content": "string"
  }
  ```

- **响应**：

  - 成功

    ```json
    {
      "comment_id": "string",
      "content": "string",
      "author_id": "string",
      "article_id": "string",
      "create_time": "2025-08-10T12:00:00Z",
      "modify_time": "2025-08-10T12:00:00Z"
    }
    ```

  - 失败

    ```json
    {
      "error": "错误信息"
    }
    ```

---

## 4. 通用说明

- 所有需要登录的接口都需在 Header 中加 `Authorization: Bearer <token>`
- 所有时间字段均为 ISO8601 格式字符串
- 所有分页接口均支持 `page` 和 `limit` 参数
- 所有响应均为 `application/json` 格式
