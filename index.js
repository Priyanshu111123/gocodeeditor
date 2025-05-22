const express = require('express');
const bodyParser = require('body-parser');
const axios = require('axios');

const app = express();
const port = 3000;

app.use(bodyParser.json());
app.use(express.static(__dirname)); // So your HTML works

// 🔐 Replace with your JDoodle credentials
const clientId = '737723caf7fa1f2b9f72e3383db4c43d';
const clientSecret = '3f8377790f7f478948adbada8444d006cf29bafe9e886601987ee541b4d7cc1c';

app.post('/compile', async (req, res) => {
    const { code, className } = req.body;

    const fullCode = code.includes('public class') ? code : `public class ${className} {\n${code}\n}`;

    try {
        const response = await axios.post('https://api.jdoodle.com/v1/execute', {
            clientId,
            clientSecret,
            script: fullCode,
            language: 'java',
            versionIndex: '4'
        });

        if (response.data.output.includes("error") || response.data.statusCode !== 200) {
            res.status(200).json({ success: false, error: response.data.output });
        } else {
            res.status(200).json({ success: true, output: response.data.output });
        }
    } catch (err) {
        res.status(500).json({ success: false, error: 'JDoodle API call failed.' });
    }
});

app.post('/run', async (req, res) => {
    const { input, code, className } = req.body;

    if (!code || !className) {
        return res.status(400).json({ output: 'Missing code or className in request.' });
    }

    const fullCode = code.includes('public class') ? code : `public class ${className} {\n${code}\n}`;

    try {
        const response = await axios.post('https://api.jdoodle.com/v1/execute', {
            clientId,
            clientSecret,
            script: fullCode,
            language: 'java',
            versionIndex: '4',
            stdin: input
        });

        res.status(200).json({ output: response.data.output });
    } catch (err) {
        console.error(err.message); // Optional: print error for debugging
        res.status(500).json({ output: 'Error running code.' });
    }
});


app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}`);
});
