package funcs

func Paginate (url *url.URL, total, size uint8) {
		pages := paginate(current, last, 3, 5, 3);
};

func paginate(current, , left, middle, right) []string {
    current = parseInt(current);
    var last = Math.ceil(total / size);
    if (last < 1 || current > last) return;
    if (current < 1) current = 1;

    var nav = '';
    var prev = current - 1;
    if (prev >= 1) {
        nav += "<a href=\"" + prefix + prev + suffix + "\" class=\"on\">上一页</a>";
    };
    /*
     * divide all page number into three parts: left, middle, right。
     * calculate the left end and right start page number。
     */
    var left_end = current - Math.ceil(middle / 2);
    var left_end2 = last - right - middle;
    if (left_end2 < left_end) left_end = left_end2;
    if (left_end < left) left_end = left;
    var right_start = last - right + 1;

    var nav_count = left + middle + right;
    if (last < nav_count) nav_count = last;
    var left_last = left;
    var right_first = left + middle + 1;
    var i, now;
    for (i = now = 1; i <= nav_count; i++) {
        if (i == left_last && now < left_end) {
            nav += "<span>...</span>";
            now = left_end + 1;
        } else if (i == right_first && now < right_start) {
            nav += "<span>...</span>";
            now = right_start + 1;
        } else {
            if (now == current) {
                nav += "<span class=\"on\">" + now + "</span>";
            }
            else {
                nav += "<a href=\"" + prefix + now + suffix + "\">" + now + "</a>";
            }
            now++;
        }
    }
    var next = current + 1
    if (next <= last) {
        nav += "<a href=\"" + prefix + next + suffix + "\" class=\"on\">下一页</a>";
    }
    return nav;
}
