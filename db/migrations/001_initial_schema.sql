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
    ('concise', 'Brief summary of key points', 'Summarize the following text concisely:'),
    ('detailed', 'Comprehensive summary with main points and supporting details', 'Provide a detailed summary of the following text:'),
    ('bullet-points', 'Key points in bullet-point format', 'Extract the main points from the following text in bullet points:'),
    ('executive', 'Business-focused summary with key insights and recommendations', 'Provide an executive summary of the following text, focusing on key business insights, implications, and actionable recommendations:');