local _M = {}

function _M.hasCoupon(phone) 

end

function _M.validatePhone(phone)
    if not phone then
        return false
    end
    if string.sub(phone, 1, 2) == "86" then
        phone = string.sub(phone, 3)
    end
    if string.sub(phone, 1, 3) == "086" then
        phone = string.sub(phone, 4)
    end
    if string.match(phone, "^1[3|5|7|8|4]%d%d%d%d%d%d%d%d%d$") then
        return true
    else 
        return false
    end
end 

function _M.validateSalt(id, salt)
    ngx.log(ngx.INFO, "id:", id, " salt:", salt)
    if not salt then
        return false
    end
    local cache = ngx.shared.scache
    local eventSalt, err = cache:get("salt:"..id)
    ngx.log(ngx.INFO, " eventSalt:", eventSalt)
    if not eventSalt or eventSalt ~= salt then
        return false
    end
    return true
end

function _M.validateEffect(id)
    local cache = ngx.shared.scache
    local effectOn, err = cache:get("eeo:"..id)
    if not effectOn then
        ngx.log(ngx.ERR, "can't get eeo ", err)
        return false
    end

    local duration, err = cache:get("ed:"..id)
    if not duration then
        ngx.log(ngx.ERR, "can't get ed ", err)
        return false
    end

    local now = ngx.now() * 1000
    if effectOn > now or now > effectOn + duration then 
        return false
    end

    return true
end

return _M
