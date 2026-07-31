"""
Menggunakan model hasil fine-tuning untuk mengklasifikasikan intent pertanyaan
sebelum di-routing ke orchestrator. Ini membuat fine-tuned model benar-benar
berperan dalam pipeline, bukan sekadar dilatih lalu tidak dipakai.
"""
from sentence_transformers import SentenceTransformer, util

MODEL_PATH = "./finetuned-support-intent-model"

_model = None

def get_model():
    global _model
    if _model is None:
        _model = SentenceTransformer(MODEL_PATH)
    return _model


INTENT_ANCHORS = {
    "retur": "bagaimana cara retur atau mengembalikan barang",
    "pengiriman": "status pengiriman dan estimasi waktu sampai",
    "pembayaran": "masalah pembayaran atau metode bayar",
}


def classify_intent(question: str) -> str:
    model = get_model()
    q_emb = model.encode(question, convert_to_tensor=True)
    best_intent, best_score = None, -1
    for intent, anchor_text in INTENT_ANCHORS.items():
        anchor_emb = model.encode(anchor_text, convert_to_tensor=True)
        score = util.cos_sim(q_emb, anchor_emb).item()
        if score > best_score:
            best_intent, best_score = intent, score
    return best_intent