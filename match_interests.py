import json
import logging
import os

os.environ["TQDM_DISABLE"] = "1"

from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

from config import HF_TOKEN

logging.getLogger("sentence_transformers").setLevel(logging.ERROR)
logging.getLogger("transformers").setLevel(logging.ERROR)


def match_people(interests: list[str]) -> list[tuple[int, int, float]]:
    model = SentenceTransformer(
        "paraphrase-multilingual-MiniLM-L12-v2",
        token=HF_TOKEN,
    )
    vectors = model.encode(interests)

    sim_matrix = cosine_similarity(vectors)

    n = len(interests)
    scores = []
    for i in range(n):
        for j in range(i + 1, n):
            scores.append((sim_matrix[i, j], i, j))

    scores.sort(reverse=True)

    used = set()
    pairs = []

    for score, i, j in scores:
        if i not in used and j not in used:
            pairs.append((i, j, round(float(score), 3)))
            used.add(i)
            used.add(j)

    return pairs


if __name__ == "__main__":
    with open("input.json") as f:
        people = json.load(f)

    pairs = match_people(people)

    output = [{"a": people[i], "b": people[j], "score": score} for i, j, score in pairs]

    with open("output.json", "w") as f:
        json.dump(output, f, ensure_ascii=False, indent=2)
