const teamupCalendarKey = 'ks1c2vhnot2ttfvawo'
const teamupToken = process.env.HUBOT_TEAMUP_TOKEN

module.exports = (robot) => {
  robot.error(function (err) {
    robot.logger.error(err)
    robot.send(`Error while running scrummaster bot: ${err}`)
  })
  robot.hear(/scrum(\s|-)?master(s)?/gi, (res) => {
    const dateString = new Date().toISOString().slice(0, 10)
    robot.http(`https://api.teamup.com/${teamupCalendarKey}/events?startDate=${dateString}&endDate=${dateString}`)
      .headers({'Teamup-Token': teamupToken})
      .get()((err, response, body) => {
        if (err) {
          robot.emit('error', `problem getting scrummaster: '${err}'`)
          return
        }
        let parsedBody = {}
        try {
          parsedBody = JSON.parse(body)
        } catch (e) {
          robot.emit('error', `problem parsing '${body}' as JSON`)
          return
        }
        if (typeof parsedBody.events === 'undefined' || parsedBody.events.length === 0) {
          robot.emit('error', `scrummasters not found: '${body}'`)
          return
        }
        const scrummaster = parsedBody.events[0].title
        res.send(`^ ${scrummaster}`)
      })
  })
}
