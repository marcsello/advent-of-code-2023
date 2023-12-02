file = io.open("input.txt")

local function trim(s)
    -- lol: http://lua-users.org/wiki/StringTrim
   return (s:gsub("^%s*(.-)%s*$", "%1"))
end

local function validate(color, num)

    if color == "red" then
        return num <= 12
    end

    if color == "green" then
        return num <= 13
    end

    if color == "blue" then
        return num <= 14
    end

end

local gameID = 1 -- Aint nobody will parse that
local sum = 0
for line in file:lines() do
    local colonI = string.find(line, ":")

    local valid = true

    for gameStr in string.gmatch(string.sub(line, colonI+1), "([^;]+)") do
        for drawStr in string.gmatch(gameStr, "([^,]+)") do
            drawStr = trim(drawStr)

            local spaceI = string.find(drawStr, " ")
            local num = string.sub(drawStr, 0, spaceI)
            local color = string.sub(drawStr, spaceI+1)

            if not validate(color, 0 + num) then
                valid = false
            end

        end
    end

    if valid then
        sum = sum + gameID
    end

    gameID = gameID + 1
end

print(sum)