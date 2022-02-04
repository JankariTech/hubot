"use strict"
const tmetricToken = process.env.HUBOT_TMETRIC_TOKEN
const tmetricAccountId = process.env.HUBOT_TMETRIC_ACCOUNT_ID
const tmetricClientList = process.env.HUBOT_TMETRIC_CLIENT_LIST
const room = 'developers'
var cron = require('node-cron');


module.exports = (robot) => {
  robot.error(function (err, res) {
    robot.logger.error(err)
    robot.send({room: room}, `there is an issue with the t-metric bot '${err}'`)
  })
  robot.hear(/tmetric\smonthly/gi, (res) => {
    tmetricMonthlyReport(
      robot,
      {
        res: res
      }
    )
  })

  robot.hear(/tmetric\sweekly/gi, (res) => {
    tmetricWeeklyReport(
      robot,
      {
        res: res
      }
    )
  })

  cron.schedule('30 12 * * 5', () => {
    tmetricWeeklyReport(
      robot,
      {
        outputPrefix: '@all ',
        outputPostfix: 'Aaphno :tea:-metric :thumbsup:-sanga :straight_ruler:-bhayo?'
      },
    )
  })

  cron.schedule('00 16 28-31 * *', () => {
    const tomorrow = new Date()
    tomorrow.setDate(new Date().getDate() + 1)
    if (tomorrow.getDate() === 1) {
      tmetricMonthlyReport(
        robot,
        {
          outputPrefix: '@all ',
          outputPostfix: 'Aaphno :tea:-metric :thumbsup:-sanga :straight_ruler:-bhayo?'
        },
      )
    }
  })
}

function tmetricMonthlyReport(robot, {outputPrefix = '', outputPostfix = '', res = null}) {
  const today = new Date()
  const startDate = `${today.getFullYear()}-${today.getMonth() + 1}-1`
  outputPrefix += `This month (since ${startDate}) we have these chargable hours logged in :tmetric:\n`
  getTmetricreport(robot, {outputPrefix, outputPostfix, startDate, res})
}

function tmetricWeeklyReport(robot, {outputPrefix = '', outputPostfix = '', res = null}) {
  const mondayOfThisWeek = getMonday(new Date())
  const startDate = `${mondayOfThisWeek.getFullYear()}-${mondayOfThisWeek.getMonth() + 1}-${mondayOfThisWeek.getDate()}`
  outputPrefix += `This week (since ${startDate}) we have these chargable hours logged in :tmetric:\n`
  getTmetricreport(robot, {outputPrefix, outputPostfix, startDate, res})
}

function getTmetricreport(robot, {outputPrefix = '', startDate, outputPostfix = '', res = null}) {
  let clientList = []
  if (tmetricToken === '' || tmetricToken === undefined) {
    robot.emit('error', `env. variable HUBOT_TMETRIC_TOKEN is not set`)
    return
  }
  if (tmetricAccountId === '' || tmetricAccountId === undefined) {
    robot.emit('error', `env. variable HUBOT_TMETRIC_ACCOUNT_ID is not set`)
    return
  }
  if (typeof tmetricClientList === 'string') {
    clientList = tmetricClientList.split(',')
  } else {
    robot.emit('error', `env. variable HUBOT_TMETRIC_CLIENT_LIST is not set or not valid`)
    return
  }

  let query = {AccountId: tmetricAccountId, StartDate: startDate, ClientList: clientList}

  robot.http(`https://app.tmetric.com/api/reports/detailed`)
    .query(query)
    .headers({Accept: 'application/json', Authorization: `Bearer ${tmetricToken}`})
    .get()((err, response, body) => {
      let parsedData = {}
      let durations = []
      let output = {}
      output.text = outputPrefix

      if (err) {
        robot.emit('error', `problem getting t-metric report: '${err}'`)
        return
      }

      try {
        parsedData = JSON.parse(body)
      } catch (e) {
        robot.emit('error', `problem parsing '${body}' as JSON`)
        return
      }

      for (const key in parsedData) {
        if (typeof durations[parsedData[key]['user']] === 'undefined') {
          durations[parsedData[key]['user']] = 0
        }
        durations[parsedData[key]['user']] += parseInt(parsedData[key]['duration'])
      }

      for (const user in durations) {
        durations[user] = parseFloat((durations[user] / 3600000).toFixed(2))
      }
      for (const user in durations) {
        output.text += user + ' ' + durations[user] + 'h\n'
      }

      output.text += outputPostfix

      if (res !== null) {
        res.send(output.text)
      } else {
        robot.send({room: room}, output.text)
      }
    })
}

function getMonday(d) {
  d = new Date(d);
  var day = d.getDay(),
    diff = d.getDate() - day + (day === 0 ? -6 : 1); // adjust when day is sunday
  return new Date(d.setDate(diff));
}
