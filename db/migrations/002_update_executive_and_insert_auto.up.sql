-- Update executive style to qna
UPDATE summarization_styles 
SET name = 'qna',
    description = 'QnA style summary',
    prompt_template = 'Transform the provided text into a QnA format by formulating pertinent questions and providing clear, direct answers based on the content.'
WHERE name = 'executive';

-- Insert new auto style
INSERT INTO summarization_styles (name, description, prompt_template) VALUES
    ('auto', 'Adaptive summary based on context', 'Summarize the following text with a length and level of detail that is optimally adapted to the content and context, ensuring the summary is both concise and comprehensive as needed.');
