# Bug 复现说明

## Bug 是什么

通知仓储的 FindByID 对不存在记录返回 (nil, nil)，service 层又忽略错误，给不存在的通知标记已读时对空指针解引用，接口 panic。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run TestFoodNilComposite
```

## 错误信息

```
panic: runtime error: invalid memory address or nil pointer dereference
github.com/blueship581/cyfreshfood/internal/service.(*NotificationService).MarkRead(...)
```
