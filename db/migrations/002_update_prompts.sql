BEGIN TRANSACTION;
-- Update prompt templates for existing summarization styles
UPDATE summarization_styles 
SET prompt_template = 'Summarize the following text concisely, focusing on the core message and key takeaways. Aim for brevity and clarity, as if writing a headline capturing the essence of the text.'
WHERE name = 'concise';

UPDATE summarization_styles 
SET prompt_template = 'Provide a detailed summary of the following text, including all main points and relevant supporting details. Ensure the summary is comprehensive and captures the nuances and context of the original text, similar to a short report.'
WHERE name = 'detailed';

UPDATE summarization_styles 
SET prompt_template = 'Extract the main points from the following text and present them as bullet points. Focus on clear, concise takeaways that are easy to scan and understand. Organize the bullet points logically for quick comprehension, like summarizing meeting minutes.'
WHERE name = 'bullet-points';

UPDATE summarization_styles 
SET prompt_template = 'Provide an executive summary of the following text, focusing on key business insights, strategic implications, and actionable recommendations.  Highlight the potential business impact, considering a leadership perspective and focusing on ROI and strategic value, as if preparing a management briefing document.'
WHERE name = 'executive';
COMMIT;