# Experiment 3: V1 Re-Run with Improved Prompts

## Hypothesis
Do V4-style prompts improve the original V1 model comparison results?

## Setup
- 6 models tested with improved "full file output" prompt
- Same task: add --model flag to queue-add.sh

## Results

| Model | V1 (original) | V1b (improved) |
|-------|--------------|----------------|
| MiniMax M2.7 | VALID_MODELS ✅ | ✅ + test file |
| Qwen3 Coder | VALID_MODELS ✅ | ✅ + test file |
| DeepSeek V3.2 | VALID_MODELS ✅ | ✅ + test file |
| Gemini Flash | VALID_MODELS ✅ | ✅ + test file |
| GPT-4.1 Mini | VALID_MODELS ✅ | ✅ + test file |
| Qwen3-30B | VALID_MODELS ✅ | ✅ (no test) |

## Key Finding
5/6 models now produce both implementation AND test files. Prompt improvement alone added test coverage across all models.
