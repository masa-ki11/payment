// 現在の年を取得
const currentYear = new Date().getFullYear();

// 開始年を設定
const startYear = 2022;

// 選択肢を動的に生成
const yearSelect = $("#year");
for (let year = startYear; year <= currentYear; year++) {
    yearSelect.append(new Option(year, year));
}

// 現在の年をデフォルトに設定
yearSelect.val(currentYear);


$(document).ready(function () {
    const monthNames = [
        "1月", "2月", "3月", "4月", "5月", "6月",
        "7月", "8月", "9月", "10月", "11月", "12月"
    ];
    const monthContainer = $("#month-container");
    const visibleMonths = ($(window).width() <= 576) ? 3 : 6; // PC画面では6つの月を表示
    const currentMonth = new Date().getMonth(); // 現在の月を取得（0から11のインデックス）
    let selectedMonthIndex = currentMonth;
    let firstVisibleMonthIndex = Math.max(0, currentMonth - Math.floor(visibleMonths / 2));

    function updateMonthDisplay() {
        monthContainer.html("");
        for (let i = 0; i < visibleMonths; i++) {
            const monthIndex = firstVisibleMonthIndex + i;
            const monthElement = $("<span>", {
                text: monthNames[monthIndex],
                "data-index": monthIndex
            }).addClass("mx-2 month-item");
            if (monthIndex === selectedMonthIndex) {
                monthElement.addClass("text-primary font-weight-bold");
            }
            monthContainer.append(monthElement);
        }
    }

    function updateVisibleMonths(delta) {
        firstVisibleMonthIndex += delta;
        if (firstVisibleMonthIndex < 0) {
            firstVisibleMonthIndex = 0;
        } else if (firstVisibleMonthIndex > monthNames.length - visibleMonths) {
            firstVisibleMonthIndex = monthNames.length - visibleMonths;
        }
        updateMonthDisplay();
    }

    updateMonthDisplay();

    $("#prev-month").on("click", function () {
        updateVisibleMonths(-1);
    });

    $("#next-month").on("click", function () {
        updateVisibleMonths(1);
    });

    monthContainer.on("click", ".month-item", function () {
        selectedMonthIndex = parseInt($(this).data("index"));
        updateMonthDisplay();
    });

    // 履歴を取得して表示する関数
    function fetchAndUpdateHistory() {
        const selectedYear = $("#year").val();
        const selectedMonth = selectedMonthIndex + 1; // 月は1から始まるため、インデックスに1を追加

        // 選択された年と月に基づいて履歴を取得するAPIを呼び出し
        const userId = $("#user-id").data("user-id");
        $.get("/user-history?year=" + selectedYear + "&month=" + selectedMonth + "&user_id=" + userId, function (data) {
            // 取得したデータを元に履歴を表示する処理を実装
            const historyTableBody = $("#history-table tbody");
            historyTableBody.html(""); // 既存の履歴をクリア

            data.forEach(function (record) {
                const newRow = $("<tr>");
                newRow.append($("<td>").text(record.created_at_formatted));
                newRow.append($("<td>").text(record.point.toLocaleString()));
                newRow.append($("<td>").text(record.action));
                newRow.append($("<td>").text(record.details));
                historyTableBody.append(newRow);
            });
        });
    }
    $("#year").on("change", fetchAndUpdateHistory);
    monthContainer.on("click", ".month-item", fetchAndUpdateHistory);

    // 初期状態で履歴を表示する
    fetchAndUpdateHistory();

});


