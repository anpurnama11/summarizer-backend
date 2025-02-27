CREATE TABLE history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    title TEXT,
    content TEXT NOT NULL,
    summary TEXT NOT NULL,
    style_id INTEGER,
    language TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (style_id) REFERENCES summarization_styles(id)
);

CREATE TABLE summarization_styles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    prompt_template TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default summarization styles
INSERT INTO summarization_styles (name, description, prompt_template) VALUES
    ('concise', 'Brief summary of key points', 'Summarize the following text concisely, focusing on the core message and key takeaways. Aim for brevity and clarity, as if writing a headline capturing the essence of the text.'),
    ('detailed', 'Comprehensive summary with main points and supporting details', 'Provide a detailed summary of the following text, including all main points and relevant supporting details. Ensure the summary is comprehensive and captures the nuances and context of the original text, similar to a short report.'),
    ('bullet-points', 'Key points in bullet-point format', 'Extract the main points from the following text and present them as bullet points. Focus on clear, concise takeaways that are easy to scan and understand. Organize the bullet points logically for quick comprehension, like summarizing meeting minutes.'),
    ('executive', 'Business-focused summary with key insights and recommendations', 'Provide an executive summary of the following text, focusing on key business insights, strategic implications, and actionable recommendations.  Highlight the potential business impact, considering a leadership perspective and focusing on ROI and strategic value, as if preparing a management briefing document.');