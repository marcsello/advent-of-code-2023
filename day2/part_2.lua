file = io.open("input.txt")

local function trim(s)
    -- lol: http://lua-users.org/wiki/StringTrim
   return (s:gsub("^%s*(.-)%s*$", "%1"))
end

local counter = {
    red= 0,
    green= 0,
    blue= 0
}

function counter:reset()
    self.red = 0
    self.green = 0
    self.blue = 0
end

function counter:consider(color, num)
    if self[color] < num then
        self[color] = num
    end
end

function counter:getPow()
    return self.red * self.green * self.blue
end

local gameID = 1 -- Aint nobody will parse that
local sum = 0
for line in file:lines() do
    local colonI = string.find(line, ":")

    counter:reset()
    for gameStr in string.gmatch(string.sub(line, colonI+1), "([^;]+)") do
        for drawStr in string.gmatch(gameStr, "([^,]+)") do
            drawStr = trim(drawStr)
            local spaceI = string.find(drawStr, " ")
            local num = string.sub(drawStr, 0, spaceI)
            local color = string.sub(drawStr, spaceI+1)
            counter:consider(color, 0 + num)
        end
    end
    sum = sum + counter:getPow()

    gameID = gameID + 1
end

print(sum)