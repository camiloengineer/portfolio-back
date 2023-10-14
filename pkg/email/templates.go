package email

const DeveloperBody = `
<p>Name: {{.Name}},</p>
<p>Email: {{.Email}},</p>
<p>Message: {{.Message}}</p>
`

const UserBody = `
<p>Hello {{.Name}},</p>
<p>I hope you are well. I have received your email and will respond to it as soon as possible.</p>
<p>Best regards.</p>
`
