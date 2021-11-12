"use strict"
const githubUsername = process.env.HUBOT_GITHUB_USERNAME
const githubToken = process.env.HUBOT_GITHUB_TOKEN
const organisation = 'owncloud'
const projectName = 'Current: QA/CI/TestAutomation'
const fastLaneColumnName = 'Fastlane'
const intervalInS = 3600
const room = 'developers'

const interval = intervalInS * 1000
const auth = Buffer.from(`${githubUsername}:${githubToken}`, 'binary').toString('base64')
module.exports = robot =>
  setInterval(() => {
    robot.http(`https://api.github.com/orgs/${organisation}/projects`)
      .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
      .get()((err, response, body) => {
          const data = JSON.parse(body).find(function (project) {
            return project.name === projectName
          })
          const columnsUrl = data.columns_url
          robot.http(columnsUrl)
            .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
            .get()((err, response, body) => {
              const data = JSON.parse(body).find(function (column) {
                return column.name === fastLaneColumnName
              })
              robot.http(data.cards_url)
                .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
                .get()((err, response, body) => {
                  const cards = JSON.parse(body)
                  let text = ''
                  if (cards.length === 1) {
                    text = `is one card`
                  } else if (cards.length > 1){
                    text = `are ${cards.length} cards`
                  } else {
                    return
                  }
                  robot.send({room: room}, `@scrummaster there ${text} in the fastlane`)
                  for (var i = 0; i < cards.length; i++) {
                    robot.http(cards[i].content_url)
                      .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
                      .get()((err, response, body) => {
                        robot.send({room: room}, JSON.parse(body).html_url)
                      })
                  }
                })
            })
        }
      );

  }, interval);
