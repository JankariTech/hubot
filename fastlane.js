"use strict"
const githubUsername = process.env.HUBOT_GITHUB_USERNAME
const githubToken = process.env.HUBOT_GITHUB_TOKEN
const organisation = 'owncloud'
const fastLaneColumnName = 'Fastlane'
const intervalInS = 3600
const room = 'developers'
const teamupCalendarKey = 'ks1c2vhnot2ttfvawo'
const teamupToken = process.env.HUBOT_TEAMUP_TOKEN

const interval = intervalInS * 1000
const auth = Buffer.from(`${githubUsername}:${githubToken}`, 'binary').toString('base64')
const GITHUB_GRAPHQL_API_URL = 'https://api.github.com/graphql';

let endCursor = null
let hasNextPage = true
let fastlaneCards = []

const getFastlaneCards = (response)=> {
    const items = response.data.data.organization.projectV2.items;
    const issueList = [];

    node.nodes.forEach(value => {
        const title = value.content.title;
        const url = value.content.url;
        const nameNode = value.fieldValues.nodes.find(field => field.name);
        const name = nameNode.name

        if (title && url && name && name === `${fastLaneColumnName}`) {
            issueList.push([title, url])
        }
    })

    return [node,issueList]
}

const generateQuery = (itemArg) => {
    return `query {
    organization(login: "${organisation}") {
    projectV2(number: 386) {
      items(${itemArg}) {
        totalCount
        nodes {
          content {
            ... on Issue {
              title
              url
            }
            ... on PullRequest {
              title
              url
            }
          }
          fieldValues(last: 3) {
            nodes {
              ... on ProjectV2ItemFieldSingleSelectValue {
                name
              }
              ... on ProjectV2ItemFieldTextValue {
                text
              }
            }
          }
        }
        pageInfo {
          endCursor
          hasNextPage
        }
      }
    }
  }
}`
}

module.exports = robot =>
    setInterval(() => {
        robot.error(function (err, res) {
            robot.logger.error(err)
            robot.send({room: room}, `there is an issue with the fastlane bot '${err}'`)
        })
        const sendRequest = () => {
            if (hasNextPage) {
                robot.http(GITHUB_GRAPHQL_API_URL)
                    .headers({Accept: 'application/json', Authorization: `Basic ${auth}`})
                    .data({query: generateQuery(`first: 30, after: "${endCursor}"`)})
                    .post()((err, response) => {
                        if (err) {
                            robot.emit('error', `problem getting projects list: '${err}'`)
                            return
                        }
                        const node = getFastlaneCards(response)[0]
                        hasNextPage = node.pageInfo.hasNextPage
                        endCursor = node.pageInfo.endCursor
                        fastlaneCards = fastlaneCards.concat(getFastlaneCards(response)[1])
                        sendRequest()
                    })
            }
        }
        sendRequest()
        let text = ''
        if (fastlaneCards.length === 1) {
            text = `is one card`
        } else if (fastlaneCards.length > 1) {
            text = `are ${fastlaneCards.length} cards`
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
                fastlaneCards.forEach(items=>{
                    robot.send({room: room}, items[1])
                })
            })

}, interval);
