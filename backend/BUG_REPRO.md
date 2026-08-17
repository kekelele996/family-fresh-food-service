# Bug 复现说明

## Bug 是什么

四个仓储 FindByID 方法使用命名返回值加 defer 把错误重置为 nil，查询不存在的食品/成员/通知/用户时静默返回成功。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run TestFoodDeferSwallow
```

## 错误信息

```
--- FAIL: TestFoodDeferSwallow
    food FindByID err=<nil>, want ErrNotFound
    member FindByID err=<nil>, want ErrNotFound
    notification FindByID err=<nil>, want ErrNotFound
    user FindByID err=<nil>, want ErrNotFound
```
