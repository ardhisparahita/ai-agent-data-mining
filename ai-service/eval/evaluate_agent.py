from rag.retriever import answer_with_rag
import time

def evaluate_efficiency(question: str, collection_name: str = "faq") -> float:
    """Efficiency: waktu respons dalam detik (makin kecil makin baik)."""
    start = time.time()
    answer_with_rag(question, collection_name)
    return round(time.time() - start, 2)


TEST_CASES = [
    {
        "question": "Berapa lama saya bisa retur barang?",
        "expected_keywords": ["7 hari"],
    },
    {
        "question": "Berapa estimasi pengiriman ke luar Jawa?",
        "expected_keywords": ["5-7 hari", "luar jawa"],
    },
]


def evaluate_effectiveness(answer: str, expected_keywords: list[str]) -> float:
    """Effectiveness: proporsi keyword kunci yang berhasil disebut dalam jawaban."""
    answer_lower = answer.lower()
    hits = sum(1 for kw in expected_keywords if kw.lower() in answer_lower)
    return hits / len(expected_keywords)


def evaluate_explainability(source_documents: list) -> bool:
    """Explainability: apakah sistem menyertakan dokumen sumber (bukan black-box)."""
    return len(source_documents) > 0


def evaluate_hallucination_risk(answer: str, source_documents: list) -> bool:
    """Hallucination check (heuristik kasar): apakah ada kata kunci jawaban
    yang tidak muncul sama sekali di dokumen sumber manapun."""
    context_text = " ".join([doc.page_content.lower() for doc in source_documents])
    answer_words = set(w.strip(".,") for w in answer.lower().split() if len(w) > 5)
    unsupported = [w for w in answer_words if w not in context_text]
    # Toleransi: kalau >50% kata "penting" tidak ada di konteks, tandai berisiko halusinasi
    return len(unsupported) > 0 and (len(unsupported) / max(len(answer_words), 1)) > 0.5
efficiency = evaluate_efficiency(case["question"])
print(f"Efficiency (waktu respons): {efficiency} detik")

def run_evaluation():
    results = []
    for case in TEST_CASES:
        result = answer_with_rag(case["question"], collection_name="faq")
        answer = result["result"]
        sources = result["source_documents"]

        effectiveness = evaluate_effectiveness(answer, case["expected_keywords"])
        explainability = evaluate_explainability(sources)
        hallucination_risk = evaluate_hallucination_risk(answer, sources)

        results.append({
            "question": case["question"],
            "answer": answer,
            "effectiveness_score": effectiveness,
            "explainability": explainability,
            "hallucination_risk": hallucination_risk,
        })

    print("=== Hasil Evaluasi Model ===\n")
    for r in results:
        print(f"Pertanyaan       : {r['question']}")
        print(f"Jawaban          : {r['answer']}")
        print(f"Effectiveness    : {r['effectiveness_score']:.2f} (0-1, makin tinggi makin baik)")
        print(f"Explainability   : {'Ya' if r['explainability'] else 'Tidak'} (ada dokumen sumber?)")
        print(f"Risiko Halusinasi: {'TINGGI' if r['hallucination_risk'] else 'Rendah'}")
        print("-" * 60)

    avg_effectiveness = sum(r["effectiveness_score"] for r in results) / len(results)
    print(f"\nRata-rata Effectiveness: {avg_effectiveness:.2f}")

    return results


if __name__ == "__main__":
    run_evaluation()