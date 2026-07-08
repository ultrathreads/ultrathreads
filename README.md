# UltraThreads

![Screenshot](screenshot.png)

`UltraThreads` is an open-source, lightweight web forum powered by Go and React, with native threaded post views.

1

## Features

- data storage in a MySQL or SQLite database
- nodes
- three possible views of posts

## Project Structure

### api

> Developed with `Go`, providing RESTful-style APIs.

*Tech Stack*
- gin (https://github.com/gin-gonic/gin) Go web framework
- JWT (https://github.com/appleboy/gin-jwt) JWT Middleware for Gin framework
- gorm (http://gorm.io/) ORM framework for Go language

### web

> Frontend page rendering service implemented based on `Next.js`.

*Tech Stack*
- Next.js (https://nextjs.org) The React Framework for the Web

## Installation

The project is still in early development. Refer to the README files of the `api` and `web` modules for installation and development instructions.

## The Principle of Threaded Posts

    📌 [ID:1] Go (ParentId:0, ThreadId:1)
        ├── 💬 [ID:11] Ordinary flat comment (ParentId:1, ThreadId:1)
        ├── 💬 [ID:12] Single nested comment (ParentId:1, ThreadId:1)
        │   └── 💬 [ID:13] One-level reply (ParentId:12, ThreadId:1)
        └── 💬 [ID:14] Multi-branch root comment (ParentId:1, ThreadId:1)
            ├── 💬 [ID:15] Branch reply A (ParentId:14, ThreadId:1)
            └── 💬 [ID:16] Branch reply B (ParentId:14, ThreadId:1)

    📌 [ID:2] Rust (ParentId:0, ThreadId:2)
        ├── 💬 [ID:21] Ultra deep nesting demo (ParentId:2, ThreadId:2)
        │   └── 💬 [ID:22] Level 1 (ParentId:21, ThreadId:2)
        │       └── 💬 [ID:23] Level 2 (ParentId:22, ThreadId:2)
        │           └── 💬 [ID:24] Level 3 (ParentId:23, ThreadId:2)
        │               └── 💬 [ID:25] Level 4 deepest (ParentId:24, ThreadId:2)
        ├── 💬 [ID:26] No child comment (ParentId:2, ThreadId:2)
        └── 💬 [ID:27] Normal reply (ParentId:2, ThreadId:2)

    📌 [ID:3] Python (ParentId:0, ThreadId:3)
        ├── 💬 [ID:31] Complex multi-branch tree (ParentId:3, ThreadId:3)
        │   ├── 💬 [ID:32] Sub branch 1 (ParentId:31, ThreadId:3)
        │   │   └── 💬 [ID:33] Sub-sub branch (ParentId:32, ThreadId:3)
        │   └── 💬 [ID:34] Sub branch 2 (ParentId:31, ThreadId:3)
        │       ├── 💬 [ID:35] Independent leaf A (ParentId:34, ThreadId:3)
        │       └── 💬 [ID:36] Independent leaf B (ParentId:34, ThreadId:3)
        ├── 💬 [ID:37] Pure flat list comment 1 (ParentId:3, ThreadId:3)
        ├── 💬 [ID:38] Pure flat list comment 2 (ParentId:3, ThreadId:3)
        └── 💬 [ID:39] Pure flat list comment 3 (ParentId:3, ThreadId:3)

    📌 [ID:4] Java (ParentId:0, ThreadId:4)
        ├── 💬 [ID:41] Alternate nested structure (ParentId:4, ThreadId:4)
        │   └── 💬 [ID:42] First layer (ParentId:41, ThreadId:4)
        │       ├── 💬 [ID:43] Side branch (ParentId:42, ThreadId:4)
        │       └── 💬 [ID:44] Main continue nested (ParentId:42, ThreadId:4)
        │           └── 💬 [ID:45] Final deepest leaf (ParentId:44, ThreadId:4)
        └── 💬 [ID:46] Only one single root child (ParentId:4, ThreadId:4)
    
## License
UltraThreads is open-sourced software licensed under the [GNU General Public License version 3](https://opensource.org/license/gpl-3.0)

