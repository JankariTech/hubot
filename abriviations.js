module.exports = (robot) => {
    robot.hear(/((\s|^)[A-Z]{3,5}(\s|$))/, (res) => {
        res.reply(`Sorry I don't understand, what does "${res.match[1]}" mean?`)
    })
}

