-- Remove auto style
DELETE FROM summarization_styles WHERE name = 'auto';

-- Revert qna style back to executive
UPDATE summarization_styles 
SET name = 'executive',
    description = 'Business-focused summary with key insights and recommendations',
    prompt_template = 'Provide an executive summary of the following text, focusing on key business insights, strategic implications, and actionable recommendations.  Highlight the potential business impact, considering a leadership perspective and focusing on ROI and strategic value, as if preparing a management briefing document.'
WHERE name = 'qna';
