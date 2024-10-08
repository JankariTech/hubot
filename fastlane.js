"use strict"
const githubUsername = process.env.HUBOT_GITHUB_USERNAME
const githubToken = process.env.HUBOT_GITHUB_TOKEN
const organisation = "owncloud"
const projectName = "QA/CI/TestAutomation"
const projectNumber = 386
const fastLaneColumnName = "Fastlane"
const itemsPerRequest = 50 // number of items (issues/PRs) to fetch from the project board
const intervalInS = 3600
const room = "developers"
const teamupCalendarKey = "ks1c2vhnot2ttfvawo"
const teamupToken = process.env.HUBOT_TEAMUP_TOKEN
const GITHUB_GRAPHQL_API_URL = "https://api.github.com/graphql"

const interval = intervalInS * 1000
const auth = Buffer.from(`${githubUsername}:${githubToken}`, "binary").toString(
  "base64"
)

let endCursor = null // graphql pagination
const fastlaneCards = []

const filterFastlaneCards = (cards) => {
  cards.forEach((card) => {
    const url = card?.content?.url
    const columnName = card?.fieldValueByName?.name

    if (url && columnName && columnName === fastLaneColumnName) {
      fastlaneCards.push(url)
    }
  })

  return fastlaneCards
}

const generateQuery = (itemArg) => {
  return `query {
    organization(login: "${organisation}"){
      projectV2(number: ${projectNumber}){
        title
        items(${itemArg}) {
          nodes{
            content{
              ...on Issue {
                url
              }
              ...on PullRequest {
                url
              }
            }
            fieldValueByName(name: "Status") {
              ...on ProjectV2ItemFieldSingleSelectValue {
                name
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

const reportFastlaneCards = (robot) =>
  robot
    .http(GITHUB_GRAPHQL_API_URL)
    .headers({ Accept: "application/json", Authorization: `Basic ${auth}` })
    .post(
      JSON.stringify({
        query: generateQuery(
          `first: ${itemsPerRequest}, after: "${endCursor}"`
        ),
      })
    )((err, _, body) => {
    if (err) {
      robot.emit("error", `problem getting projects list: '${err}'`)
      return
    }
    const parsedBody = JSON.parse(body)
    
    if (Object.hasOwn(parsedBody, "errors")) {
      const errorTypes = parsedBody.errors.map((error) => error.type)
      if (errorTypes.every((type) => type === "INSUFFICIENT_SCOPES")) {
        robot.emit(
          "error",
          `could not find project '${projectName}' . Do you have the right permissions?`
        )
        return
      }

      if (errorTypes.every((type) => type === "NOT_FOUND")) {
        console.log("I am error")
        robot.emit(
          "error",
          `could not find project'${projectName}' with number ${projectNumber}`
        )
        return
      }
    }

    if (!parsedBody.data) {
      robot.emit("error", `Response doesn't have data: '${err}'`)
      return
    }

    if (parsedBody.data.organization.projectV2.title !== projectName) {
      robot.emit("error", `Project title from response doesn't match`)
      return
    }

    const items = parsedBody.data.organization.projectV2.items
    endCursor = items.pageInfo.endCursor

    filterFastlaneCards(items.nodes)

    if (items.pageInfo.hasNextPage) {
      reportFastlaneCards(robot)
    } else {
      let text = ""
      if (fastlaneCards.length === 1) {
        text = `is one card`
      } else if (fastlaneCards.length > 1) {
        text = `are ${fastlaneCards.length} cards`
      } else {
        return
      }
      pingScrumMasters(robot, text)
    }
  })

const pingScrumMasters = (robot, text) => {
  const dateString = new Date().toISOString().slice(0, 10)
  robot
    .http(
      `https://api.teamup.com/${teamupCalendarKey}/events?startDate=${dateString}&endDate=${dateString}`
    )
    .headers({ "Teamup-Token": teamupToken })
    .get()((err, response, body) => {
    if (err) {
      robot.emit("error", `problem getting scrummaster: '${err}'`)
      return
    }
    let parsedBody = {}
    try {
      parsedBody = JSON.parse(body)
    } catch (e) {
      robot.emit("error", `problem parsing '${body}' as JSON`)
      return
    }
    if (
      typeof parsedBody.events === "undefined" ||
      parsedBody.events.length === 0
    ) {
      robot.emit("error", `scrummasters not found: '${body}'`)
      return
    }
    const scrummaster = parsedBody.events[0].title
    // NOTE: Link preview is not shown if there are more than 5 links.
    const links = fastlaneCards.join("\n")
    robot.send(
      { room },
      `${scrummaster} there ${text} in the fastlane\n${links}`
    )
  })
}

module.exports = (robot) => {
  robot.error(function (err) {
    robot.logger.error(err)
    robot.send({ room: room }, `Error while running fastlane bot: ${err}`)
  })

  return setInterval(() => {
    endCursor = null
    fastlaneCards.length = 0

    reportFastlaneCards(robot)
  }, interval)
}