"use strict"
const githubUsername = process.env.HUBOT_GITHUB_USERNAME
const githubToken = process.env.HUBOT_GITHUB_TOKEN
const organisation = 'owncloud'
const projectName = 'Current: QA/CI/TestAutomation'
const fastLaneColumnName = 'Fastlane'
const intervalInS = 3600
const room = 'developers'
const teamupCalendarKey = 'ks1c2vhnot2ttfvawo'
const teamupToken = process.env.HUBOT_TEAMUP_TOKEN

const interval = intervalInS * 1000
const auth = Buffer.from(`${githubUsername}:${githubToken}`, 'binary').toString('base64')

module.exports = robot =>
  setInterval(() => {
    robot.error(function (err, res) {
      robot.logger.error(err)
      robot.send({room: room}, `there is an issue with the fastlane bot '${err}'`)
    })
    robot.http(`https://api.github.com/orgs/${organisation}/projects`)
      .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
      .get()((err, response, body) => {
          if (err) {
            robot.emit('error', `problem getting projects list: '${err}'`)
            return
          }
          let parsedData = {}
          try {
            parsedData = JSON.parse(body)
          } catch (e) {
            robot.emit('error', `problem parsing '${body}' as JSON`)
            return
          }

          if (!Array.isArray(parsedData)) {
            robot.emit('error', `Response body cannot be parsed to an array. Content: "${body}"`)
            return
          }

          const data = parsedData.find(function (project) {
            return project.name === projectName
          })

          if (typeof data !== 'object' || typeof data.columns_url !== 'string' ) {
            robot.emit('error', `could not find project '${projectName}' do you have the right permissions?`)
            return
          }

          robot.http(data.columns_url)
            .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
            .get()((err, response, body) => {
              if (err) {
                robot.emit('error', `problem getting columns list: '${err}'`)
                return
              }
              let parsedData = {}
              try {
                parsedData = JSON.parse(body)
              } catch (e) {
                robot.emit('error', `problem parsing '${body}' as JSON`)
                return
              }
              const data = parsedData.find(function (column) {
                return column.name === fastLaneColumnName
              })
              if (typeof data !== 'object' || typeof data.cards_url !== 'string' ) {
                robot.emit('error', `could not find column '${fastLaneColumnName}'`)
                return
              }
              robot.http(data.cards_url)
                .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
                .get()((err, response, body) => {
                  if (err) {
                    robot.emit('error', `problem getting cards list: '${err}'`)
                    return
                  }
                  let cards = {}
                  try {
                    cards = JSON.parse(body)
                  } catch (e) {
                    robot.emit('error', `problem parsing '${body}' as JSON`)
                    return
                  }
                  let text = ''
                  if (cards.length === 1) {
                    text = `is one card`
                  } else if (cards.length > 1) {
                    text = `are ${cards.length} cards`
                  } else {
                    return
                  }
                  const dateString = new Date().toISOString().slice(0,10)
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
                        robot.emit('error', `problem getting scrummaster: '${body}'`)
                        return
                      }
                      const scrummaster = parsedBody.events[0].title
                      robot.send({room: room}, `${scrummaster} there ${text} in the fastlane`)
                      for (var i = 0; i < cards.length; i++) {
                        robot.http(cards[i].content_url)
                          .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
                          .get()((err, response, body) => {
                            robot.send({room: room}, JSON.parse(body).html_url)
                          })
                      }
                    })
                })
            })
        }
      );

  }, interval);
