# Bug 复现说明

## Bug 是什么

仓储层把“不存在”错误改成了普通 error，导致 service 层用 errors.Is 判断 sentinel 时全部失败，登录不存在的手机号、查询不存在的用户/家庭组/食品都返回 500 而不是 401/404。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run TestFoodErrorComposite
```

## 错误信息

```
--- FAIL: TestFoodErrorComposite
    login missing user code=..., want 1001
    get missing user code=..., want 1003
    get missing group code=..., want 1003
    get missing food code=..., want 1003
```
